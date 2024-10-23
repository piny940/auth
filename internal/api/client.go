package api

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"context"
	"errors"
)

// ClientsInterfaceGetClient implements StrictServerInterface.
func (s *Server) ClientsInterfaceGetClient(ctx context.Context, request ClientsInterfaceGetClientRequestObject) (ClientsInterfaceGetClientResponseObject, error) {
	client, err := s.ClientUsecase.Find(oauth.ClientID(request.Id))
	if errors.Is(err, domain.ErrRecordNotFound) {
		return ClientsInterfaceGetClient400JSONResponse{
			Error:            ClientNotFound,
			ErrorDescription: "client not found",
		}, nil
	}
	return ClientsInterfaceGetClient200JSONResponse{
		Client: PublicClient{
			Id:   string(client.ID),
			Name: client.Name,
		},
	}, nil
}

// AccountClientsCreateClient implements StrictServerInterface.
func (s *Server) AccountClientsCreateClient(ctx context.Context, request AccountClientsCreateClientRequestObject) (AccountClientsCreateClientResponseObject, error) {
	panic("unimplemented")
}

// AccountClientsDeleteClient implements StrictServerInterface.
func (s *Server) AccountClientsDeleteClient(ctx context.Context, request AccountClientsDeleteClientRequestObject) (AccountClientsDeleteClientResponseObject, error) {
	panic("unimplemented")
}

// AccountClientsListClients implements StrictServerInterface.
func (s *Server) AccountClientsListClients(ctx context.Context, request AccountClientsListClientsRequestObject) (AccountClientsListClientsResponseObject, error) {
	panic("unimplemented")
}

// AccountClientsUpdateClient implements StrictServerInterface.
func (s *Server) AccountClientsUpdateClient(ctx context.Context, request AccountClientsUpdateClientRequestObject) (AccountClientsUpdateClientResponseObject, error) {
	panic("unimplemented")
}
