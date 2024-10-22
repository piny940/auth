package api

import (
	"auth/internal/domain/oauth"
	"auth/internal/usecase"
	"context"
	"errors"
	"net/url"
	"strings"
)

func (s *Server) OAuthInterfaceAuthorize(ctx context.Context, request OAuthInterfaceAuthorizeRequestObject) (OAuthInterfaceAuthorizeResponseObject, error) {
	user, err := CurrentUser(ctx)
	if errors.Is(err, ErrUnauthorized) {
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
		authorizeUrl, err := url.JoinPath(s.Conf.ServerUrl, "authorize")
		if err != nil {
			return nil, err
		}
		next := authorizeUrl + "?" + toQueryString(query)
		url := s.Conf.LoginUrl + "?" + toQueryString(map[string]string{"next": next})

		return OAuthInterfaceAuthorize302Response{
			Headers: OAuthInterfaceAuthorize302ResponseHeaders{
				Location: url,
			},
		}, nil
	}
	if err != nil {
		return nil, err
	}
	err = s.AuthUsecase.Request(user, toDAuthParams(request.Params))
	if errors.Is(err, oauth.ErrInvalidRequestType) {
		return OAuthInterfaceAuthorize400JSONResponse{
			Error:            UnsupportedResponseType,
			ErrorDescription: "unsupported_response_type",
		}, nil
	}
	if errors.Is(err, usecase.ErrNotApproved) {
		// TODO: redirect to approve page
		return OAuthInterfaceAuthorize400JSONResponse{
			Error:            AccessDenied,
			ErrorDescription: "access_denied",
		}, nil
	}
	if errors.Is(err, oauth.ErrInvalidScope) {
		return OAuthInterfaceAuthorize400JSONResponse{
			Error:            InvalidScope,
			ErrorDescription: "invalid_scope",
		}, nil
	}
	if errors.Is(err, oauth.ErrInvalidRedirectURI) {
		return OAuthInterfaceAuthorize400JSONResponse{
			Error:            InvalidRequest,
			ErrorDescription: "redirect_uri is invalid",
		}, nil
	}
	if err != nil {
		return nil, err
	}
	// TODO: attach auth code
	return OAuthInterfaceAuthorize302Response{
		Headers: OAuthInterfaceAuthorize302ResponseHeaders{
			Location: request.Params.RedirectUri,
		},
	}, nil
}

// OAuthInterfaceGetToken implements StrictServerInterface.
func (s *Server) OAuthInterfaceGetToken(ctx context.Context, request OAuthInterfaceGetTokenRequestObject) (OAuthInterfaceGetTokenResponseObject, error) {
	panic("unimplemented")
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
