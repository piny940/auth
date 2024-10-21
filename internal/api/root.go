package api

import (
	"auth/internal/domain/oauth"
	"auth/internal/usecase"
	"errors"
	"net/http"
	"strings"

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
	req := toDAuthParams(params)
	user, err := CurrentUser(ctx.Request())
	if err != nil {
		return err
	}
	if user == nil {
		return ctx.Redirect(http.StatusFound, "/?error=unauthorized_client")
	}
	err = s.AuthUsecase.Request(user, req)
	if errors.Is(err, oauth.ErrInvalidRequestType{}) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":             "unsupported_response_type",
			"error_description": "unsupported_response_type",
		})
	}
	if errors.Is(err, usecase.ErrNotApproved{}) {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{
			"error":             "access_denied",
			"error_description": "access_denied",
		})
	}
	if errors.Is(err, oauth.ErrInvalidScope{}) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":             "invalid_scope",
			"error_description": "invalid_scope",
		})
	}
	if errors.Is(err, oauth.ErrInvalidRedirectURI{}) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":             "invalid_request",
			"error_description": "invalid_request",
		})
	}
	if err != nil {
		return err
	}
	return ctx.Redirect(http.StatusFound, req.RedirectURI)
}

func toDAuthParams(params AuthorizeParams) *oauth.AuthRequest {
	strScopes := strings.Split(params.Scope, " ")
	scopes := make([]oauth.TypeScope, 0, len(strScopes))
	for _, s := range strScopes {
		scopes = append(scopes, oauth.TypeScope(s))
	}
	return &oauth.AuthRequest{
		ClientID:     oauth.ClientID(params.ClientId),
		RedirectURI:  params.RedirectUri,
		ResponseType: oauth.TypeResponseType(params.ResponseType),
		Scopes:       scopes,
		State:        params.State,
	}
}
