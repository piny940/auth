package middleware

import (
	"auth/internal/api"
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
)

type AuthConfig struct {
	Issuer       string `required:"true" envconfig:"OAUTH_ISSUER"`
	RsaPublicKey string `required:"true" envconfig:"OAUTH_RSA_PUBLIC_KEY"`
}

type AuthMiddleware struct {
	clientRepo oauth.IClientRepo
	userRepo   domain.IUserRepo
	issuer     string
	rsaPubKey  *rsa.PublicKey
}

func NewAuthMiddleware(clientRepo oauth.IClientRepo, userRepo domain.IUserRepo) *AuthMiddleware {
	conf := &AuthConfig{}
	err := envconfig.Process("auth", conf)
	if err != nil {
		panic(err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(conf.RsaPublicKey))
	if err != nil {
		panic(err)
	}
	return &AuthMiddleware{
		clientRepo: clientRepo,
		userRepo:   userRepo,
		issuer:     conf.Issuer,
		rsaPubKey:  pubKey,
	}
}

// IMPORTANT: This middleware is dependent on the session middleware.
// Make sure to add the session middleware before this middleware.
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

func (m *AuthMiddleware) authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	if input.RequestValidationInput.Request.URL.Path == "/oauth/token" && input.SecuritySchemeName == "BasicAuth" {
		return m.authClient(ctx, input)
	}
	if input.SecuritySchemeName == "BearerAuth" {

	}
	if input.SecuritySchemeName != "ApiKeyAuth" || input.SecurityScheme.In != "cookie" {
		return ErrUnsupportedAuthenticationScheme
	}
	return m.cookieAuth(ctx, input)
}

func (m *AuthMiddleware) bearerAuth(c context.Context, input *openapi3filter.AuthenticationInput) error {
	auth := input.RequestValidationInput.Request.Header.Get("Authorization")
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
	api.SetScopes(c, &api.AccessScope{User: user})
	return err
}

func (m *AuthMiddleware) cookieAuth(_ context.Context, input *openapi3filter.AuthenticationInput) error {
	_, err := api.CurrentUser(input.RequestValidationInput.Request.Context())
	return err
}

func (m *AuthMiddleware) authClient(_ context.Context, input *openapi3filter.AuthenticationInput) error {
	auth := input.RequestValidationInput.Request.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		return fmt.Errorf("invalid basic auth header")
	}
	decoded, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return err
	}
	clientArr := strings.Split(string(decoded), ":")
	if len(clientArr) != 2 {
		return fmt.Errorf("invalid basic auth header")
	}
	client, err := m.clientRepo.FindByID(oauth.ClientID(clientArr[0]))
	if err != nil {
		return err
	}
	return client.SecretCorrect(clientArr[1])
}

var ErrUnsupportedAuthenticationScheme = fmt.Errorf("unsupported authentication scheme")
