package gateway

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"errors"
	"fmt"
	"slices"

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

const redirectUriBatchSize = 20

func (c *ClientRepo) FindByID(id oauth.ClientID) (*oauth.Client, error) {
	cq := c.query.Client
	rq := c.query.RedirectURI
	client, err := cq.Where(cq.ID.Eq(string(id))).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrRecordNotFound
	}
	uris, err := rq.Where(rq.ClientID.Eq(client.ID)).Find()
	if err != nil {
		return nil, err
	}
	return toDomainClient(client, uris), nil
}

func (c *ClientRepo) FindWithUserID(id oauth.ClientID, userID domain.UserID) (*oauth.Client, error) {
	fmt.Println("clientId: ", id, "userid: ", userID)
	client, err := c.query.Client.Where(
		c.query.Client.ID.Eq(string(id)), c.query.Client.UserID.Eq(int64(userID)),
	).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	uris, err := c.query.RedirectURI.Where(c.query.RedirectURI.ClientID.Eq(client.ID)).Find()
	if err != nil {
		return nil, err
	}
	return toDomainClient(client, uris), nil
}

func (c *ClientRepo) Create(client *oauth.ClientInput) error {
	err := c.query.Client.Create(&model.Client{
		ID:              string(client.ID),
		EncryptedSecret: client.EncryptedSecret,
		UserID:          int64(client.UserID),
		Name:            client.Name,
	})
	if err != nil {
		return err
	}
	redirectUris := make([]*model.RedirectURI, 0, len(client.RedirectURIs))
	for _, uri := range client.RedirectURIs {
		redirectUris = append(redirectUris, &model.RedirectURI{
			ClientID: string(client.ID),
			URI:      uri,
		})
	}
	if err := c.query.RedirectURI.CreateInBatches(redirectUris, redirectUriBatchSize); err != nil {
		return err
	}
	return nil
}

func (c *ClientRepo) Update(client *oauth.Client, userID domain.UserID) error {
	_, err := c.query.Client.Where(
		c.query.Client.ID.Eq(string(client.ID)), c.query.Client.UserID.Eq(int64(userID))).
		UpdateColumns(&model.Client{
			ID:              string(client.ID),
			EncryptedSecret: client.EncryptedSecret,
			UserID:          int64(client.UserID),
			Name:            client.Name,
		})
	if err != nil {
		return err
	}
	current, err := c.query.RedirectURI.Where(c.query.RedirectURI.ClientID.Eq(string(client.ID))).Find()
	if err != nil {
		return err
	}
	adds := make([]string, 0, len(client.RedirectURIs))
	removeIds := make([]int64, 0, len(current))
	for _, uri := range client.RedirectURIs {
		if !slices.ContainsFunc(current, func(redirectUri *model.RedirectURI) bool {
			return redirectUri.URI == uri
		}) {
			adds = append(adds, uri)
		}
	}
	for _, uri := range current {
		if !slices.Contains(client.RedirectURIs, uri.URI) {
			removeIds = append(removeIds, uri.ID)
		}
	}
	addRedirects := make([]*model.RedirectURI, 0, len(adds))
	for _, uri := range adds {
		addRedirects = append(addRedirects, &model.RedirectURI{
			ClientID: string(client.ID),
			URI:      uri,
		})
	}
	if err := c.query.RedirectURI.CreateInBatches(addRedirects, redirectUriBatchSize); err != nil {
		return err
	}
	if _, err := c.query.RedirectURI.Where(c.query.RedirectURI.ID.In(removeIds...)).Delete(); err != nil {
		return err
	}
	return nil
}

func (c *ClientRepo) Delete(id oauth.ClientID, userID domain.UserID) error {
	_, err := c.query.Client.Where(
		c.query.Client.ID.Eq(string(id)), c.query.Client.UserID.Eq(int64(userID)),
	).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientRepo) List(userID domain.UserID) ([]*oauth.Client, error) {
	clients, err := c.query.Client.Where(c.query.Client.UserID.Eq(int64(userID))).Find()
	if err != nil {
		return nil, err
	}
	clientIDs := make([]string, 0, len(clients))
	for _, client := range clients {
		clientIDs = append(clientIDs, client.ID)
	}
	uris, err := c.query.RedirectURI.Where(c.query.RedirectURI.ClientID.In(clientIDs...)).Find()
	if err != nil {
		return nil, err
	}
	urisByClientID := make(map[string][]*model.RedirectURI)
	for _, uri := range uris {
		if _, ok := urisByClientID[uri.ClientID]; !ok {
			urisByClientID[uri.ClientID] = make([]*model.RedirectURI, 0, 1)
		}
		urisByClientID[uri.ClientID] = append(urisByClientID[uri.ClientID], uri)
	}
	result := make([]*oauth.Client, 0, len(clients))
	for _, client := range clients {
		result = append(result, toDomainClient(client, urisByClientID[client.ID]))
	}
	return result, nil
}

func toDomainClient(client *model.Client, redirectUris []*model.RedirectURI) *oauth.Client {
	uris := make([]string, 0, len(redirectUris))
	for _, r := range redirectUris {
		uris = append(uris, r.URI)
	}
	return &oauth.Client{
		ID:              oauth.ClientID(client.ID),
		EncryptedSecret: client.EncryptedSecret,
		UserID:          domain.UserID(client.UserID),
		Name:            client.Name,
		CreatedAt:       client.CreatedAt,
		UpdatedAt:       client.UpdatedAt,
		RedirectURIs:    uris,
	}
}
