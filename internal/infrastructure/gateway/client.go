package gateway

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ClientRepo struct {
	db    *infrastructure.DB
	query *query.Query
}

func NewClientRepo(db *infrastructure.DB) oauth.IClientRepo {
	query := query.Use(db.Client)
	return &ClientRepo{
		db:    db,
		query: query,
	}
}

var _ oauth.IClientRepo = &ClientRepo{}

func (c *ClientRepo) FindByID(id oauth.ClientID) (*oauth.Client, error) {
	cq := c.query.Client
	rq := c.query.RedirectURI
	client, err := cq.Where(cq.ID.Eq(string(id))).First()
	fmt.Println(id, client)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrRecordNotFound
	}
	uris, err := rq.Where(rq.ClientID.Eq(client.ID)).Find()
	if err != nil {
		return nil, err
	}
	return toDomainClient(client, uris), nil
}

func toDomainClient(client *model.Client, redirectUris []*model.RedirectURI) *oauth.Client {
	uris := make([]string, 0, len(redirectUris))
	for _, r := range redirectUris {
		uris = append(uris, r.URI)
	}
	return &oauth.Client{
		ID:              oauth.ClientID(client.ID),
		EncryptedSecret: client.EncryptedSecret,
		UserID:          client.UserID,
		Name:            client.Name,
		CreatedAt:       client.CreatedAt,
		UpdatedAt:       client.UpdatedAt,
		RedirectURIs:    uris,
	}
}
