package api

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/usecase"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func (s *Server) OAuthInterfaceAuthorize(ctx context.Context, request OAuthInterfaceAuthorizeRequestObject) (OAuthInterfaceAuthorizeResponseObject, error) {
	thisUrl := func() (string, error) {
		query := map[string]string{
			"error":         "unauthorized_client",
			"redirect_uri":  request.Params.RedirectUri,
			"response_type": string(request.Params.ResponseType),
			"client_id":     request.Params.ClientId,
			"scope":         request.Params.Scope,
		}
		if request.Params.State != nil {
			query["state"] = *request.Params.State
		}
		authorizeUrl, err := url.JoinPath(s.Conf.ServerUrl, "oauth", "authorize")
		if err != nil {
			return "", fmt.Errorf("failed to join path: %w", err)
		}
		return authorizeUrl + "?" + toQueryString(query), nil
	}
	loginUrl := func() (string, error) {
		this, err := thisUrl()
		if err != nil {
			return "", fmt.Errorf("failed to get this url: %w", err)
		}
		query := map[string]string{
			"next":              this,
			"client_id":         request.Params.ClientId,
			"scope":             request.Params.Scope,
			"response_type":     string(request.Params.ResponseType),
			"error":             string(OAuthAuthorizeErrUnauthorizedClient),
			"error_description": "unauthorized_client",
		}
		if request.Params.State != nil {
			query["state"] = *request.Params.State
		}
		return s.Conf.LoginUrl + "?" + toQueryString(query), nil
	}

	authSession, err := s.Auth.CurrentUser(ctx)
	if errors.Is(err, ErrUnauthorized) {
		loginUrl, err := loginUrl()
		if err != nil {
			return nil, err
		}
		return OAuthInterfaceAuthorize302Response{
			Headers: OAuthInterfaceAuthorize302ResponseHeaders{
				Location: loginUrl,
			},
		}, nil
	}
	if err != nil {
		return nil, err
	}
	session := &oauth.Session{
		User:     authSession.User,
		AuthTime: authSession.AuthTime,
	}
	if request.Params.ResponseType != Code {
		s.logger.Infof("invalid request type: %v", err)
		return OAuthInterfaceAuthorize400JSONResponse{
			Error:            OAuthAuthorizeErrUnsupportedResponseType,
			ErrorDescription: "unsupported_response_type",
		}, nil
	}
	code, err := s.OAuthUsecase.RequestCodeAuth(session, toUAuthParams(request.Params))
	if errors.Is(err, oauth.ErrSessionExpired) {
		s.logger.Infof("session expired: %v", err)
		loginUrl, err := loginUrl()
		if err != nil {
			return nil, err
		}
		return OAuthInterfaceAuthorize302Response{
			Headers: OAuthInterfaceAuthorize302ResponseHeaders{
				Location: loginUrl,
			},
		}, nil
	}
	if errors.Is(err, usecase.ErrClientNotFound) {
		s.logger.Infof("invalid client id: %v", err)
		return OAuthInterfaceAuthorize400JSONResponse{
			Error:            OAuthAuthorizeErrInvalidRequest,
			ErrorDescription: err.Error(),
		}, nil
	}
	if errors.Is(err, usecase.ErrNotApproved) {
		this, err := thisUrl()
		if err != nil {
			return nil, err
		}
		query := map[string]string{
			"next":              this,
			"client_id":         request.Params.ClientId,
			"scope":             request.Params.Scope,
			"error":             string(OAuthAuthorizeErrAccessDenied),
			"error_description": "access_denied",
		}
		if request.Params.State != nil {
			query["state"] = *request.Params.State
		}
		url := s.Conf.ApproveUrl + "?" + toQueryString(query)
		return OAuthInterfaceAuthorize302Response{
			Headers: OAuthInterfaceAuthorize302ResponseHeaders{
				Location: url,
			},
		}, nil
	}
	if errors.Is(err, oauth.ErrInvalidScope) {
		s.logger.Infof("invalid scope: %v", err)
		return OAuthInterfaceAuthorize400JSONResponse{
			Error:            OAuthAuthorizeErrInvalidScope,
			ErrorDescription: "invalid_scope",
		}, nil
	}
	if errors.Is(err, usecase.ErrRedirectURINotRegistered) {
		s.logger.Infof("invalid redirect uri: %v", err)
		return OAuthInterfaceAuthorize400JSONResponse{
			Error:            OAuthAuthorizeErrInvalidRequest,
			ErrorDescription: "redirect_uri is invalid",
		}, nil
	}
	if err != nil {
		s.logger.Errorf("failed to authorize: %v", err)
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
		GrantType:   string(request.Body.GrantType),
		AuthCode:    request.Body.Code,
		RedirectURI: request.Body.RedirectUri,
		ClientID:    oauth.ClientID(request.Body.ClientId),
	})
	if errors.Is(err, usecase.ErrInvalidGrantType) {
		s.logger.Infof("invalid grant type: %v", err)
		return OAuthInterfaceGetToken400JSONResponse{
			Error:            InvalidRequest,
			ErrorDescription: "invalid grant type",
		}, nil
	}
	if errors.Is(err, oauth.ErrInvalidRedirectURI) {
		s.logger.Infof("invalid redirect uri: %v", err)
		return OAuthInterfaceGetToken400JSONResponse{
			Error:            InvalidRequest,
			ErrorDescription: "invalid redirect uri",
		}, nil
	}
	if errors.Is(err, domain.ErrRecordNotFound) ||
		errors.Is(err, oauth.ErrInvalidClientID) ||
		errors.Is(err, oauth.ErrExpiredAuthCode) ||
		errors.Is(err, oauth.ErrUsedAuthCode) {
		s.logger.Infof("invalid auth code: %v", err)
		return OAuthInterfaceGetToken400JSONResponse{
			Error:            InvalidRequest,
			ErrorDescription: "invalid auth code",
		}, nil
	}
	if err != nil {
		s.logger.Errorf("failed to get token: %v", err)
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

func toUAuthParams(params OAuthInterfaceAuthorizeParams) *usecase.AuthRequest {
	strScopes := strings.Split(params.Scope, " ")
	scopes := make([]oauth.TypeScope, 0, len(strScopes))
	for _, s := range strScopes {
		scopes = append(scopes, oauth.TypeScope(s))
	}
	var maxAge *time.Duration
	if params.MaxAge != nil {
		d := time.Duration(*params.MaxAge) * time.Second
		maxAge = &d
	}
	return &usecase.AuthRequest{
		ClientID:     oauth.ClientID(params.ClientId),
		RedirectURI:  params.RedirectUri,
		ResponseType: usecase.TypeResponseType(params.ResponseType),
		Scopes:       scopes,
		State:        params.State,
		MaxAge:       maxAge,
	}
}

func (s *Server) OAuthInterfaceGetJwks(ctx context.Context, request OAuthInterfaceGetJwksRequestObject) (OAuthInterfaceGetJwksResponseObject, error) {
	set, err := s.OAuthUsecase.GetJWKs()
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(set)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return OAuthInterfaceGetJwks200JSONResponse(res), nil
}
