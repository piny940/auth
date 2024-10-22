package api

import (
	"auth/internal/usecase"

	"github.com/kelseyhightower/envconfig"
)

type Server struct {
	AuthUsecase *usecase.AuthUsecase
	Conf        *Config
}

type Config struct {
	ServerUrl  string `envconfig:"SERVER_URL" required:"true"`
	LoginUrl   string `split_words:"true" required:"true"`
	ApproveUrl string `split_words:"true" required:"true"`
}

var _ StrictServerInterface = &Server{}

func NewServer(authUsecase *usecase.AuthUsecase) *Server {
	conf := &Config{}
	err := envconfig.Process("api", conf)
	if err != nil {
		panic(err)
	}
	return &Server{
		AuthUsecase: authUsecase,
		Conf:        conf,
	}
}

// // Authorize implements ServerInterface.
// func (s *Server) Authorize(ctx echo.Context, params AuthorizeParams) error {
// 	req := toDAuthParams(params)
// 	user, err := CurrentUser(ctx.Request())
// 	if err != nil {
// 		return err
// 	}
// 	if user == nil {
// 		query := map[string]string{
// 			"error":         "unauthorized_client",
// 			"redirect_uri":  params.RedirectUri,
// 			"response_type": params.ResponseType,
// 			"client_id":     params.ClientId,
// 			"scope":         params.Scope,
// 		}
// 		if params.State != nil {
// 			query["state"] = *params.State
// 		}
// 		authorizeUrl, err := url.JoinPath(s.Conf.ServerUrl, "authorize")
// 		if err != nil {
// 			return err
// 		}
// 		next := authorizeUrl + "?" + toQueryString(query)
// 		url := s.Conf.LoginUrl + "?" + toQueryString(map[string]string{"next": next})
// 		return ctx.Redirect(http.StatusFound, url)
// 	}
// 	err = s.AuthUsecase.Request(user, req)
// 	if errors.Is(err, oauth.ErrInvalidRequestType) {
// 		return ctx.JSON(http.StatusBadRequest, echo.Map{
// 			"error":             "unsupported_response_type",
// 			"error_description": "unsupported_response_type",
// 		})
// 	}
// 	if errors.Is(err, usecase.ErrNotApproved) {
// 		return ctx.JSON(http.StatusUnauthorized, echo.Map{
// 			"error":             "access_denied",
// 			"error_description": "access_denied",
// 		})
// 	}
// 	if errors.Is(err, oauth.ErrInvalidScope) {
// 		return ctx.JSON(http.StatusBadRequest, echo.Map{
// 			"error":             "invalid_scope",
// 			"error_description": "invalid_scope",
// 		})
// 	}
// 	if errors.Is(err, oauth.ErrInvalidRedirectURI) {
// 		return ctx.JSON(http.StatusBadRequest, echo.Map{
// 			"error":             "invalid_request",
// 			"error_description": "redirect_uri is invalid",
// 		})
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return ctx.Redirect(http.StatusFound, req.RedirectURI)
// }

// // PostAuthorize implements ServerInterface.
// func (s *Server) PostAuthorize(ctx echo.Context) error {
// 	panic("unimplemented")
// }

// func toDAuthParams(params AuthorizeParams) *oauth.AuthRequest {
// 	strScopes := strings.Split(params.Scope, " ")
// 	scopes := make([]oauth.TypeScope, 0, len(strScopes))
// 	for _, s := range strScopes {
// 		scopes = append(scopes, oauth.TypeScope(s))
// 	}
// 	return &oauth.AuthRequest{
// 		ClientID:     oauth.ClientID(params.ClientId),
// 		RedirectURI:  params.RedirectUri,
// 		ResponseType: oauth.TypeResponseType(params.ResponseType),
// 		Scopes:       scopes,
// 		State:        params.State,
// 	}
// }

// func toQueryString(h map[string]string) string {
// 	var q []string
// 	for k, v := range h {
// 		q = append(q, k+"="+url.QueryEscape(v))
// 	}
// 	return strings.Join(q, "&")
// }
