package api

import (
	"auth/internal/domain"
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
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":             "invalid_name_or_password",
			"error_description": "name or password is incorrect",
		})
	}
	err = Login(ctx.Request(), ctx.Response().Writer, user)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// Signup implements ServerInterface.
func (s *Server) Signup(ctx echo.Context) error {
	var body SignupJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	err := s.AuthUsecase.SignUp(body.Name, body.Password, body.PasswordConfirmation)
	if errors.Is(err, domain.ErrNameLengthNotEnough{}) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":             "name_length_not_enough",
			"error_description": err.Error(),
		})
	}
	if errors.Is(err, domain.ErrNameAlreadyUsed{}) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":             "name_already_used",
			"error_description": err.Error(),
		})
	}
	if errors.Is(err, domain.ErrPasswordLengthNotEnough{}) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":             "password_length_not_enough",
			"error_description": err.Error(),
		})
	}
	if errors.Is(err, domain.ErrPasswordConfirmation{}) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":             "password_confirmation",
			"error_description": err.Error(),
		})
	}
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusNoContent, nil)
}

// Me implements ServerInterface.
func (s *Server) Me(ctx echo.Context) error {
	user, err := CurrentUser(ctx.Request())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
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

// PostAuthorize implements ServerInterface.
func (s *Server) PostAuthorize(ctx echo.Context) error {
	panic("unimplemented")
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
