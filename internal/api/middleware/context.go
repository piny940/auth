package middleware

import (
	"auth/internal/api"
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
)

type EchoContextMiddleware struct{}

func NewEchoContextMiddleware() *EchoContextMiddleware {
	return &EchoContextMiddleware{}
}

type typeEchoContextKey string

const echoContextKey typeEchoContextKey = "echoContext"

func (m *EchoContextMiddleware) Context() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			newReq := req.WithContext(context.WithValue(req.Context(), echoContextKey, c))

			*req = *newReq
			return next(c)
		}
	}
}

type EchoContextReg struct{}

var _ api.EchoContextReg = &EchoContextReg{}

func NewEchoContextReg() *EchoContextReg {
	return &EchoContextReg{}
}

func (r *EchoContextReg) Context(c context.Context) (echo.Context, error) {
	echoContext, ok := c.Value(echoContextKey).(echo.Context)
	if !ok {
		return nil, fmt.Errorf("failed to get echo.Context from context")
	}
	return echoContext, nil
}
