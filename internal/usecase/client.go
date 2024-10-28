package usecase

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"fmt"
)

type ClientUsecase struct {
	ClientRepo oauth.IClientRepo
}

type IClientUsecase interface {
	Find(id oauth.ClientID) (*oauth.Client, error)
	List(userID domain.UserID) ([]*oauth.Client, error)
	Create(userID domain.UserID, name string, redirectURIs []string) (*CreatedClient, error)
	Delete(id oauth.ClientID, userID domain.UserID) error
}

func NewClientUsecase(clientRepo oauth.IClientRepo) IClientUsecase {
	return &ClientUsecase{
		ClientRepo: clientRepo,
	}
}

var _ IClientUsecase = &ClientUsecase{}

func (c *ClientUsecase) Find(id oauth.ClientID) (*oauth.Client, error) {
	return c.ClientRepo.FindByID(id)
}

type CreatedClient struct {
	ID           oauth.ClientID
	Name         string
	Secret       string
	RedirectURIs []string
}

func (c *ClientUsecase) Create(userID domain.UserID, name string, redirectURIs []string) (*CreatedClient, error) {
	clientID, err := oauth.IssueClientID()
	if err != nil {
		return nil, fmt.Errorf("failed to issue client id: %w", err)
	}
	rawSecret, encryptedSecret, err := oauth.IssueClientSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to issue client secret: %w", err)
	}
	if err := c.ClientRepo.Create(&oauth.ClientInput{
		ID:              clientID,
		Name:            name,
		EncryptedSecret: encryptedSecret,
		UserID:          userID,
		RedirectURIs:    redirectURIs,
	}); err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &CreatedClient{ID: clientID, Name: name, Secret: rawSecret, RedirectURIs: redirectURIs}, nil
}

func (c *ClientUsecase) Delete(id oauth.ClientID, userID domain.UserID) error {
	if err := c.ClientRepo.Delete(id, userID); err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}
	return nil
}

func (c *ClientUsecase) List(userID domain.UserID) ([]*oauth.Client, error) {
	clients, err := c.ClientRepo.List(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list clients: %w", err)
	}
	return clients, nil
}
