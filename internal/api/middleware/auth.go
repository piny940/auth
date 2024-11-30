package middleware

import (
	"auth/internal/api"
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
)

type AuthMiddleware struct {
	clientRepo   oauth.IClientRepo
	userRepo     domain.IUserRepo
	issuer       string
	rsaPubKey    *rsa.PublicKey
	sessionStore *sessions.CookieStore
}

const sessionName = "auth.middleware"
const sessionUserKey = "user"
const sessionAuthTimeKey = "auth_time"

type typeAuthContextKey string

const userContextKey typeAuthContextKey = "current_user"
const scopesContextKey typeAuthContextKey = "scopes"

func NewAuthMiddleware(conf *Config, clientRepo oauth.IClientRepo, userRepo domain.IUserRepo) *AuthMiddleware {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(conf.RsaPublicKey))
	if err != nil {
		panic(err)
	}
	store := sessions.NewCookieStore([]byte(conf.SessionSecret))
	return &AuthMiddleware{
		clientRepo:   clientRepo,
		userRepo:     userRepo,
		issuer:       conf.Issuer,
		rsaPubKey:    pubKey,
		sessionStore: store,
	}
}

func (m *AuthMiddleware) Auth() echo.MiddlewareFunc {
	spec, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}
	spec.Servers = nil // HACK: https://github.com/deepmap/oapi-codegen/blob/master/examples/petstore-expanded/echo/petstore.go#L30-L32
	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: m.authenticate,
			},
		})
	return validator
}

var ErrUnsupportedAuthenticationScheme = fmt.Errorf("unsupported authentication scheme")

func (m *AuthMiddleware) authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	if input.SecuritySchemeName == "BasicAuth" {
		return m.authClient(ctx, input)
	}
	if input.SecuritySchemeName == "BearerAuth" {
		return m.bearerAuth(ctx, input)
	}
	if input.SecuritySchemeName != "ApiKeyAuth" || input.SecurityScheme.In != "cookie" {
		return ErrUnsupportedAuthenticationScheme
	}
	return m.cookieAuth(ctx, input)
}

func (m *AuthMiddleware) bearerAuth(_ context.Context, input *openapi3filter.AuthenticationInput) error {
	req := input.RequestValidationInput.Request
	auth := req.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return fmt.Errorf("invalid bearer token")
	}
	token := auth[7:]
	tok, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return m.rsaPubKey, nil
	})
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}
	claims := tok.Claims.(jwt.MapClaims)
	if claims["iss"] != m.issuer {
		return fmt.Errorf("invalid issuer")
	}
	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		return fmt.Errorf("token expired")
	}
	sub := claims["sub"].(string)
	idStr := strings.Split(strings.Split(sub, ";")[0], ":")[1]
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid user id")
	}
	user, err := m.userRepo.FindByID(domain.UserID(userID))
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	authTime := int64(claims["auth_time"].(float64))
	authSession := &api.AuthSession{
		User:     user,
		AuthTime: time.Unix(authTime, 0),
	}
	scopesStr := strings.Split(claims["scope"].(string), " ")
	scopes := make([]oauth.TypeScope, 0, len(scopesStr))
	for _, s := range scopesStr {
		scopes = append(scopes, oauth.TypeScope(s))
	}
	newCtx := context.WithValue(req.Context(), userContextKey, authSession)
	newCtx = context.WithValue(newCtx, scopesContextKey, scopes)
	newReq := req.WithContext(newCtx)
	*req = *newReq
	return err
}

func (m *AuthMiddleware) cookieAuth(_ context.Context, input *openapi3filter.AuthenticationInput) error {
	req := input.RequestValidationInput.Request
	authSession, err := authSessionFromCookie(req, m.sessionStore)
	if err != nil {
		return fmt.Errorf("failed to get auth session from cookie: %w", err)
	}
	newReq := req.WithContext(context.WithValue(req.Context(), userContextKey, authSession))
	*req = *newReq
	return nil
}

func (m *AuthMiddleware) authClient(_ context.Context, input *openapi3filter.AuthenticationInput) error {
	auth := input.RequestValidationInput.Request.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		return fmt.Errorf("invalid basic auth header. expected: Basic <base64>")
	}
	decoded, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return err
	}
	clientArr := strings.Split(string(decoded), ":")
	if len(clientArr) != 2 {
		return fmt.Errorf("invalid basic auth header. expected: Basic <client_id>:<client_secret>")
	}
	client, err := m.clientRepo.FindByID(oauth.ClientID(clientArr[0]))
	if err != nil {
		return err
	}
	return client.SecretCorrect(clientArr[1])
}

type Auth struct {
	echoContextReg *EchoContextReg
	store          *sessions.CookieStore
}

var _ api.Auth = &Auth{}

func NewAuth(conf *Config, echoContextReg *EchoContextReg) api.Auth {
	return &Auth{
		echoContextReg: echoContextReg,
		store:          sessions.NewCookieStore([]byte(conf.SessionSecret)),
	}
}

func (a *Auth) CurrentUser(ctx context.Context) (*api.AuthSession, error) {
	authSession, ok := ctx.Value(userContextKey).(*api.AuthSession)
	if ok {
		return authSession, nil
	}
	echoCtx, err := a.echoContextReg.Context(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get echo context: %w", err)
	}
	req := echoCtx.Request()
	authSession, err = authSessionFromCookie(req, a.store)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth session from cookie: %w", err)
	}
	return authSession, nil
}

func (a *Auth) AccessScopes(ctx context.Context) ([]oauth.TypeScope, error) {
	scopes, ok := ctx.Value(scopesContextKey).([]oauth.TypeScope)
	if !ok {
		return nil, api.ErrUnauthorized
	}
	return scopes, nil
}

func (a *Auth) Login(ctx context.Context, user *domain.User) error {
	echoCtx, err := a.echoContextReg.Context(ctx)
	if err != nil {
		return fmt.Errorf("failed to get echo context: %w", err)
	}
	reg := sessions.GetRegistry(echoCtx.Request())
	session, err := reg.Get(a.store, sessionName)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	session.Values[sessionUserKey] = user
	session.Values[sessionAuthTimeKey] = time.Now().Unix()
	if err := session.Save(echoCtx.Request(), echoCtx.Response().Writer); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

func (a *Auth) Logout(ctx context.Context) error {
	echoCtx, err := a.echoContextReg.Context(ctx)
	if err != nil {
		return fmt.Errorf("failed to get echo context: %w", err)
	}
	reg := sessions.GetRegistry(echoCtx.Request())
	session, err := reg.Get(a.store, sessionName)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	delete(session.Values, sessionUserKey)
	delete(session.Values, sessionAuthTimeKey)
	if err := session.Save(echoCtx.Request(), echoCtx.Response().Writer); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

func authSessionFromCookie(req *http.Request, store *sessions.CookieStore) (*api.AuthSession, error) {
	reg := sessions.GetRegistry(req)
	session, err := reg.Get(store, sessionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	userObj, ok := session.Values[sessionUserKey]
	if !ok {
		return nil, api.ErrUnauthorized
	}
	user, ok := userObj.(*domain.User)
	if !ok {
		return nil, fmt.Errorf("failed to get user from session")
	}
	authTime, ok := session.Values[sessionAuthTimeKey]
	if !ok {
		return nil, api.ErrUnauthorized
	}
	t, ok := authTime.(int64)
	if !ok {
		return nil, fmt.Errorf("failed to get auth time from session")
	}
	return &api.AuthSession{
		User:     user,
		AuthTime: time.Unix(t, 0),
	}, nil
}
