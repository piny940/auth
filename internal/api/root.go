package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// V1Login implements ServerInterface.
func (s *Server) Login(ctx echo.Context) error {
	var body LoginJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}

	user, err := s.AuthUsecase.Login(body.Name, body.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid username or password")
	}
	err = Login(ctx.Request(), ctx.Response().Writer, user)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// Authorize implements ServerInterface.
func (s *Server) Authorize(ctx echo.Context, params AuthorizeParams) error {
	panic("unimplemented")
}
