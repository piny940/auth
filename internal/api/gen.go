// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// Client defines model for Client.
type Client struct {
	Id           int64    `json:"id"`
	Name         string   `json:"name"`
	RedirectUrls []string `json:"redirect_urls"`
}

// ClientCreate defines model for ClientCreate.
type ClientCreate struct {
	Name         string   `json:"name"`
	RedirectUrls []string `json:"redirect_urls"`
}

// ReqLogin defines model for ReqLogin.
type ReqLogin struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// User defines model for User.
type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// UserCreate defines model for UserCreate.
type UserCreate struct {
	Name                 string `json:"name"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// ClientsListClientsParams defines parameters for ClientsListClients.
type ClientsListClientsParams struct {
	Cookie string `json:"cookie"`
}

// ClientsCreateClientJSONBody defines parameters for ClientsCreateClient.
type ClientsCreateClientJSONBody struct {
	Client ClientCreate `json:"client"`
}

// ClientsUpdateClientJSONBody defines parameters for ClientsUpdateClient.
type ClientsUpdateClientJSONBody struct {
	Client ClientCreate `json:"client"`
}

// AuthorizeParams defines parameters for Authorize.
type AuthorizeParams struct {
	ResponseType string  `form:"response_type" json:"response_type"`
	ClientId     string  `form:"client_id" json:"client_id"`
	RedirectUri  string  `form:"redirect_uri" json:"redirect_uri"`
	Scope        string  `form:"scope" json:"scope"`
	State        *string `form:"state,omitempty" json:"state,omitempty"`
}

// PostAuthorizeMultipartBody defines parameters for PostAuthorize.
type PostAuthorizeMultipartBody struct {
	ClientId     string  `json:"client_id"`
	RedirectUri  string  `json:"redirect_uri"`
	ResponseType string  `json:"response_type"`
	Scope        string  `json:"scope"`
	State        *string `json:"state,omitempty"`
}

// TokenGetTokenJSONBody defines parameters for TokenGetToken.
type TokenGetTokenJSONBody struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// ClientsCreateClientJSONRequestBody defines body for ClientsCreateClient for application/json ContentType.
type ClientsCreateClientJSONRequestBody ClientsCreateClientJSONBody

// ClientsUpdateClientJSONRequestBody defines body for ClientsUpdateClient for application/json ContentType.
type ClientsUpdateClientJSONRequestBody ClientsUpdateClientJSONBody

// PostAuthorizeMultipartRequestBody defines body for PostAuthorize for multipart/form-data ContentType.
type PostAuthorizeMultipartRequestBody PostAuthorizeMultipartBody

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = ReqLogin

// SignupJSONRequestBody defines body for Signup for application/json ContentType.
type SignupJSONRequestBody = UserCreate

// TokenGetTokenJSONRequestBody defines body for TokenGetToken for application/json ContentType.
type TokenGetTokenJSONRequestBody TokenGetTokenJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all clients
	// (GET /account/clients/)
	ClientsListClients(ctx echo.Context, params ClientsListClientsParams) error
	// Create a new client
	// (POST /account/clients/)
	ClientsCreateClient(ctx echo.Context) error
	// Delete a client
	// (DELETE /account/clients/:id/{id})
	ClientsDeleteClient(ctx echo.Context, id int64) error
	// Update a client
	// (POST /account/clients/:id/{id})
	ClientsUpdateClient(ctx echo.Context, id int64) error
	// Authorization Request
	// (GET /authorize)
	Authorize(ctx echo.Context, params AuthorizeParams) error
	// Authorization Request
	// (POST /authorize)
	PostAuthorize(ctx echo.Context) error
	// Get me
	// (GET /me)
	Me(ctx echo.Context) error
	// Logout
	// (DELETE /session)
	Logout(ctx echo.Context) error
	// Login
	// (POST /session)
	Login(ctx echo.Context) error
	// Signup
	// (POST /signup)
	Signup(ctx echo.Context) error
	// Get a token
	// (POST /token/)
	TokenGetToken(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ClientsListClients converts echo context to params.
func (w *ServerInterfaceWrapper) ClientsListClients(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params ClientsListClientsParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "cookie" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("cookie")]; found {
		var Cookie string
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for cookie, got %d", n))
		}

		err = runtime.BindStyledParameterWithOptions("simple", "cookie", valueList[0], &Cookie, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter cookie: %s", err))
		}

		params.Cookie = Cookie
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter cookie is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ClientsListClients(ctx, params)
	return err
}

// ClientsCreateClient converts echo context to params.
func (w *ServerInterfaceWrapper) ClientsCreateClient(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ClientsCreateClient(ctx)
	return err
}

// ClientsDeleteClient converts echo context to params.
func (w *ServerInterfaceWrapper) ClientsDeleteClient(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ClientsDeleteClient(ctx, id)
	return err
}

// ClientsUpdateClient converts echo context to params.
func (w *ServerInterfaceWrapper) ClientsUpdateClient(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ClientsUpdateClient(ctx, id)
	return err
}

// Authorize converts echo context to params.
func (w *ServerInterfaceWrapper) Authorize(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params AuthorizeParams
	// ------------- Required query parameter "response_type" -------------

	err = runtime.BindQueryParameter("form", false, true, "response_type", ctx.QueryParams(), &params.ResponseType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter response_type: %s", err))
	}

	// ------------- Required query parameter "client_id" -------------

	err = runtime.BindQueryParameter("form", false, true, "client_id", ctx.QueryParams(), &params.ClientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter client_id: %s", err))
	}

	// ------------- Required query parameter "redirect_uri" -------------

	err = runtime.BindQueryParameter("form", false, true, "redirect_uri", ctx.QueryParams(), &params.RedirectUri)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter redirect_uri: %s", err))
	}

	// ------------- Required query parameter "scope" -------------

	err = runtime.BindQueryParameter("form", false, true, "scope", ctx.QueryParams(), &params.Scope)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scope: %s", err))
	}

	// ------------- Optional query parameter "state" -------------

	err = runtime.BindQueryParameter("form", false, false, "state", ctx.QueryParams(), &params.State)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter state: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Authorize(ctx, params)
	return err
}

// PostAuthorize converts echo context to params.
func (w *ServerInterfaceWrapper) PostAuthorize(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostAuthorize(ctx)
	return err
}

// Me converts echo context to params.
func (w *ServerInterfaceWrapper) Me(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Me(ctx)
	return err
}

// Logout converts echo context to params.
func (w *ServerInterfaceWrapper) Logout(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Logout(ctx)
	return err
}

// Login converts echo context to params.
func (w *ServerInterfaceWrapper) Login(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Login(ctx)
	return err
}

// Signup converts echo context to params.
func (w *ServerInterfaceWrapper) Signup(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Signup(ctx)
	return err
}

// TokenGetToken converts echo context to params.
func (w *ServerInterfaceWrapper) TokenGetToken(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.TokenGetToken(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/account/clients/", wrapper.ClientsListClients)
	router.POST(baseURL+"/account/clients/", wrapper.ClientsCreateClient)
	router.DELETE(baseURL+"/account/clients/:id/:id", wrapper.ClientsDeleteClient)
	router.POST(baseURL+"/account/clients/:id/:id", wrapper.ClientsUpdateClient)
	router.GET(baseURL+"/authorize", wrapper.Authorize)
	router.POST(baseURL+"/authorize", wrapper.PostAuthorize)
	router.GET(baseURL+"/me", wrapper.Me)
	router.DELETE(baseURL+"/session", wrapper.Logout)
	router.POST(baseURL+"/session", wrapper.Login)
	router.POST(baseURL+"/signup", wrapper.Signup)
	router.POST(baseURL+"/token/", wrapper.TokenGetToken)

}
