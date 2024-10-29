package api

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"context"
	"errors"
)

// public api
func (s *Server) ClientsInterfaceGetClient(ctx context.Context, request ClientsInterfaceGetClientRequestObject) (ClientsInterfaceGetClientResponseObject, error) {
	client, err := s.ClientUsecase.Find(oauth.ClientID(request.Id))
	if errors.Is(err, domain.ErrRecordNotFound) {
		s.logger.Infof("client not found: %v", err)
		return ClientsInterfaceGetClient400JSONResponse{
			Error:            ClientNotFound,
			ErrorDescription: "client not found",
		}, nil
	}
	if err != nil {
		s.logger.Errorf("failed to get client: %v", err)
		return nil, err
	}
	return ClientsInterfaceGetClient200JSONResponse{
		Client: PublicClient{
			Id:   string(client.ID),
			Name: client.Name,
		},
	}, nil
}

// ----------------------------------- private api -------------------------------

func (s *Server) AccountClientsCreateClient(ctx context.Context, request AccountClientsCreateClientRequestObject) (AccountClientsCreateClientResponseObject, error) {
	user, err := CurrentUser(ctx)
	if err != nil {
		s.logger.Errorf("failed to get current user: %v", err)
		return nil, err
	}
	client, err := s.ClientUsecase.Create(
		domain.UserID(user.ID), request.Body.Client.Name, request.Body.Client.RedirectUrls,
	)
	if err != nil {
		s.logger.Errorf("failed to create client: %v", err)
		return nil, err
	}
	return AccountClientsCreateClient201JSONResponse{
		Client: AccountClientsCreatedClient{
			Id:           string(client.ID),
			Name:         client.Name,
			RedirectUrls: client.RedirectURIs,
			Secret:       client.Secret,
		},
	}, nil
}

func (s *Server) AccountClientsDeleteClient(ctx context.Context, request AccountClientsDeleteClientRequestObject) (AccountClientsDeleteClientResponseObject, error) {
	user, err := CurrentUser(ctx)
	if err != nil {
		s.logger.Errorf("failed to get current user: %v", err)
		return nil, err
	}
	if err := s.ClientUsecase.Delete(oauth.ClientID(request.Id), domain.UserID(user.ID)); err != nil {
		s.logger.Errorf("failed to delete client: %v", err)
		return nil, err
	}
	return AccountClientsDeleteClient204Response{}, nil
}

func (s *Server) AccountClientsListClients(ctx context.Context, request AccountClientsListClientsRequestObject) (AccountClientsListClientsResponseObject, error) {
	user, err := CurrentUser(ctx)
	if err != nil {
		s.logger.Errorf("failed to get current user: %v", err)
		return nil, err
	}
	clients, err := s.ClientUsecase.List(domain.UserID(user.ID))
	if err != nil {
		s.logger.Errorf("failed to list clients: %v", err)
		return nil, err
	}
	mClients := make([]Client, 0, len(clients))
	for _, client := range clients {
		mClients = append(mClients, Client{
			Id:           string(client.ID),
			Name:         client.Name,
			RedirectUrls: client.RedirectURIs,
		})
	}
	return AccountClientsListClients200JSONResponse{
		Clients: mClients,
	}, nil
}

func (s *Server) AccountClientsGetClient(ctx context.Context, request AccountClientsGetClientRequestObject) (AccountClientsGetClientResponseObject, error) {
	user, err := CurrentUser(ctx)
	if err != nil {
		s.logger.Errorf("failed to get current user: %v", err)
		return nil, err
	}
	client, err := s.ClientUsecase.FindWithUserID(oauth.ClientID(request.Id), user.ID)
	if err != nil {
		s.logger.Errorf("failed to get client: %v", err)
		return AccountClientsGetClient404JSONResponse{
			Error: "client_not_found",
		}, nil
	}
	return AccountClientsGetClient200JSONResponse{
		Client: Client{
			Id:           string(client.ID),
			Name:         client.Name,
			RedirectUrls: client.RedirectURIs,
		},
	}, nil
}

func (s *Server) AccountClientsUpdateClient(ctx context.Context, request AccountClientsUpdateClientRequestObject) (AccountClientsUpdateClientResponseObject, error) {
	user, err := CurrentUser(ctx)
	if err != nil {
		s.logger.Errorf("failed to get current user: %v", err)
		return nil, err
	}
	if err := s.ClientUsecase.Update(user.ID, request.Body.Client.Name, request.Body.Client.RedirectUrls); err != nil {
		s.logger.Errorf("failed to update client: %v", err)
		return nil, err
	}
	return AccountClientsUpdateClient204Response{}, nil
}
