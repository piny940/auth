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
	FindWithUserID(id oauth.ClientID, userID domain.UserID) (*oauth.Client, error)
	List(userID domain.UserID) ([]*oauth.Client, error)
	Create(userID domain.UserID, name string, redirectURIs []string) (*CreatedClient, error)
	Update(clientID oauth.ClientID, userID domain.UserID, name string, redirectURIs []string) error
	Delete(id oauth.ClientID, userID domain.UserID) error
}

func NewClientUsecase(clientRepo oauth.IClientRepo) IClientUsecase {
	return &ClientUsecase{
		ClientRepo: clientRepo,
	}
}

var _ IClientUsecase = &ClientUsecase{}

type CreatedClient struct {
	ID           oauth.ClientID
	Name         string
	Secret       string
	RedirectURIs []string
}

func (c *ClientUsecase) Find(id oauth.ClientID) (*oauth.Client, error) {
	return c.ClientRepo.FindByID(id)
}
func (c *ClientUsecase) FindWithUserID(id oauth.ClientID, userID domain.UserID) (*oauth.Client, error) {
	fmt.Println("api: clientid: ", id)
	client, err := c.ClientRepo.FindWithUserID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find client: %w", err)
	}
	return client, nil
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

func (c *ClientUsecase) Update(clientID oauth.ClientID, userID domain.UserID, name string, redirectURIs []string) error {
	current, err := c.ClientRepo.FindWithUserID(oauth.ClientID(clientID), userID)
	if err != nil {
		return fmt.Errorf("failed to find client: %w", err)
	}
	if err := c.ClientRepo.Update(&oauth.ClientInput{
		ID:              current.ID,
		EncryptedSecret: current.EncryptedSecret,
		UserID:          current.UserID,
		Name:            name,
		RedirectURIs:    redirectURIs,
	}, userID); err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}
	return nil
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
