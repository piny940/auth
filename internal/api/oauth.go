package api

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/usecase"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"
)

func (s *Server) OAuthInterfaceAuthorize(ctx context.Context, request OAuthInterfaceAuthorizeRequestObject) (OAuthInterfaceAuthorizeResponseObject, error) {
	var this string
	{
		query := map[string]string{
			"error":         "unauthorized_client",
			"redirect_uri":  request.Params.RedirectUri,
			"response_type": request.Params.ResponseType,
			"client_id":     request.Params.ClientId,
			"scope":         request.Params.Scope,
		}
		if request.Params.State != nil {
			query["state"] = *request.Params.State
		}
		authorizeUrl, err := url.JoinPath(s.Conf.ServerUrl, "oauth", "authorize")
		if err != nil {
			return nil, err
		}
		this = authorizeUrl + "?" + toQueryString(query)
	}
	user, err := CurrentUser(ctx)
	if errors.Is(err, ErrUnauthorized) {
		url := s.Conf.LoginUrl + "?" + toQueryString(map[string]string{"next": this})
		return OAuthInterfaceAuthorize302Response{
			Headers: OAuthInterfaceAuthorize302ResponseHeaders{
				Location: url,
			},
		}, nil
	}
	if err != nil {
		return nil, err
	}
	code, err := s.OAuthUsecase.RequestAuthorization(user, toDAuthParams(request.Params))
	if errors.Is(err, oauth.ErrInvalidRequestType) {
		return OAuthInterfaceAuthorize400JSONResponse{
			Error:            OAuthAuthorizeErrUnsupportedResponseType,
			ErrorDescription: "unsupported_response_type",
		}, nil
	}
	if errors.Is(err, oauth.ErrInvalidClientID) {
		return OAuthInterfaceAuthorize400JSONResponse{
			Error:            OAuthAuthorizeErrInvalidRequest,
			ErrorDescription: err.Error(),
		}, nil
	}
	if errors.Is(err, usecase.ErrNotApproved) {
		url := s.Conf.ApproveUrl + "?" + toQueryString(map[string]string{
			"next":              this,
			"client_id":         request.Params.ClientId,
			"scope":             request.Params.Scope,
			"error":             string(OAuthAuthorizeErrAccessDenied),
			"error_description": "access_denied",
		})
		return OAuthInterfaceAuthorize302Response{
			Headers: OAuthInterfaceAuthorize302ResponseHeaders{
				Location: url,
			},
		}, nil
	}
	if errors.Is(err, oauth.ErrInvalidScope) {
		return OAuthInterfaceAuthorize400JSONResponse{
			Error:            OAuthAuthorizeErrInvalidScope,
			ErrorDescription: "invalid_scope",
		}, nil
	}
	if errors.Is(err, oauth.ErrInvalidRedirectURI) {
		return OAuthInterfaceAuthorize400JSONResponse{
			Error:            OAuthAuthorizeErrInvalidRequest,
			ErrorDescription: "redirect_uri is invalid",
		}, nil
	}
	if err != nil {
		return nil, err
	}
	query := map[string]string{
		"code": code.Value,
	}
	if request.Params.State != nil {
		query["state"] = *request.Params.State
	}
	url := request.Params.RedirectUri + "?" + toQueryString(query)
	return OAuthInterfaceAuthorize302Response{
		Headers: OAuthInterfaceAuthorize302ResponseHeaders{
			Location: url,
		},
	}, nil
}

// OAuthInterfaceGetToken implements StrictServerInterface.
func (s *Server) OAuthInterfaceGetToken(ctx context.Context, request OAuthInterfaceGetTokenRequestObject) (OAuthInterfaceGetTokenResponseObject, error) {
	accessToken, idToken, err := s.OAuthUsecase.RequestToken(&usecase.TokenRequest{
		GrantType:   request.Body.GrantType,
		AuthCode:    request.Body.Code,
		RedirectURI: request.Body.RedirectUri,
		ClientID:    oauth.ClientID(request.Body.ClientId),
	})
	if errors.Is(err, usecase.ErrInvalidGrantType) {
		return OAuthInterfaceGetToken400JSONResponse{
			Error:            InvalidRequest,
			ErrorDescription: "invalid grant type",
		}, nil
	}
	if errors.Is(err, oauth.ErrInvalidRedirectURI) {
		return OAuthInterfaceGetToken400JSONResponse{
			Error:            InvalidRequest,
			ErrorDescription: "invalid redirect uri",
		}, nil
	}
	if errors.Is(err, domain.ErrRecordNotFound) ||
		errors.Is(err, oauth.ErrInvalidClientID) ||
		errors.Is(err, oauth.ErrExpiredAuthCode) ||
		errors.Is(err, oauth.ErrUsedAuthCode) {
		slog.Info(fmt.Sprintf("invalid auth code: %s", err.Error()))
		return OAuthInterfaceGetToken400JSONResponse{
			Error:            InvalidRequest,
			ErrorDescription: "invalid auth code",
		}, nil
	}
	if err != nil {
		return nil, err
	}
	var idTokenStr *string
	if idToken != nil {
		idTokenStr = &idToken.Value
	}
	return OAuthInterfaceGetToken200JSONResponse{
		Body: OAuthTokenRes{
			AccessToken: accessToken.Value,
			IdToken:     idTokenStr,
			TokenType:   Bearer,
			ExpiresIn:   int32(time.Until(accessToken.ExpiresAt).Seconds()),
		},
	}, nil
}

// OAuthInterfacePostAuthorize implements StrictServerInterface.
func (s *Server) OAuthInterfacePostAuthorize(ctx context.Context, request OAuthInterfacePostAuthorizeRequestObject) (OAuthInterfacePostAuthorizeResponseObject, error) {
	panic("unimplemented")
}

func toDAuthParams(params OAuthInterfaceAuthorizeParams) *oauth.AuthRequest {
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
