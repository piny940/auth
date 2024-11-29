package middleware

import (
	"auth/internal/api"

	"github.com/labstack/echo/v4"
)

func Session() echo.MiddlewareFunc {
	return SetCurrentUser
}

func SetCurrentUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		api.SetSessionRegistry(c.Request())
		return next(c)
	}
}
