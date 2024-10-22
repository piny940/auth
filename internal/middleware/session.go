package middleware

import (
	"auth/internal/api"

	"github.com/labstack/echo/v4"
)

func SetCurrentUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		api.SetCurrentUserToCtx(c)
		return next(c)
	}
}
