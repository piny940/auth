package middleware

import (
	"auth/internal/api"
	"auth/internal/domain/oauth"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
)

type AuthMiddleware struct {
	clientRepo oauth.IClientRepo
}

func NewAuthMiddleware(clientRepo oauth.IClientRepo) *AuthMiddleware {
	return &AuthMiddleware{
		clientRepo: clientRepo,
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
	if input.SecuritySchemeName != "ApiKeyAuth" || input.SecurityScheme.In != "cookie" {
		return ErrUnsupportedAuthenticationScheme
	}
	return m.cookieAuth(ctx, input)
}

func (m *AuthMiddleware) cookieAuth(ctx context.Context, _ *openapi3filter.AuthenticationInput) error {
	_, err := api.CurrentUser(ctx)
	return err
}

func (m *AuthMiddleware) authClient(_ context.Context, input *openapi3filter.AuthenticationInput) error {
	raw := input.RequestValidationInput.Request.Header.Get("Authorization")[6:]
	decoded, err := base64.StdEncoding.DecodeString(raw)
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
