package api

import "context"

// PostAuthorize implements StrictServerInterface.
func (s *Server) PostAuthorize(ctx context.Context, request PostAuthorizeRequestObject) (PostAuthorizeResponseObject, error) {
	panic("unimplemented")
}

// TokenGetToken implements StrictServerInterface.
func (s *Server) TokenGetToken(ctx context.Context, request TokenGetTokenRequestObject) (TokenGetTokenResponseObject, error) {
	panic("unimplemented")
}

// Authorize implements StrictServerInterface.
func (s *Server) Authorize(ctx context.Context, request AuthorizeRequestObject) (AuthorizeResponseObject, error) {
	panic("unimplemented")
}
