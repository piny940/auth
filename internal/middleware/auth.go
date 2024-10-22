package middleware

import (
	"auth/internal/api"
	"context"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
)

func CreateAuthMiddleware() (echo.MiddlewareFunc, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, err
	}
	spec.Servers = nil // HACK: https://github.com/deepmap/oapi-codegen/blob/master/examples/petstore-expanded/echo/petstore.go#L30-L32
	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: Authenticate,
			},
		})
	return validator, nil
}

func Authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	_, err := api.UserFromReq(input.RequestValidationInput.Request)
	if err != nil {
		return err
	}
	return nil
}
