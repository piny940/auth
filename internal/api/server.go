package api

import (
	"auth/internal/usecase"
	"context"
)

type Server struct {
	AuthUsecase usecase.AuthUsecase
}

var _ StrictServerInterface = &Server{}

func NewServer(authUsecase usecase.AuthUsecase) *Server {
	return &Server{
		AuthUsecase: authUsecase,
	}
}

// ClientsCreateClient implements StrictServerInterface.
func (s *Server) ClientsCreateClient(ctx context.Context, request ClientsCreateClientRequestObject) (ClientsCreateClientResponseObject, error) {
	panic("unimplemented")
}

// ClientsDeleteClient implements StrictServerInterface.
func (s *Server) ClientsDeleteClient(ctx context.Context, request ClientsDeleteClientRequestObject) (ClientsDeleteClientResponseObject, error) {
	panic("unimplemented")
}

// ClientsListClients implements StrictServerInterface.
func (s *Server) ClientsListClients(ctx context.Context, request ClientsListClientsRequestObject) (ClientsListClientsResponseObject, error) {
	panic("unimplemented")
}

// ClientsUpdateClient implements StrictServerInterface.
func (s *Server) ClientsUpdateClient(ctx context.Context, request ClientsUpdateClientRequestObject) (ClientsUpdateClientResponseObject, error) {
	panic("unimplemented")
}

// TokenGetToken implements StrictServerInterface.
func (s *Server) TokenGetToken(ctx context.Context, request TokenGetTokenRequestObject) (TokenGetTokenResponseObject, error) {
	panic("unimplemented")
}

// V1Authorize implements StrictServerInterface.
func (s *Server) V1Authorize(ctx context.Context, request V1AuthorizeRequestObject) (V1AuthorizeResponseObject, error) {
	panic("unimplemented")
}
