package gateway

import (
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/query"
	"time"
)

type ClientRepo struct {
	db    *infrastructure.DB
	query *query.Query
}

type ClientResult struct {
	ID              string
	EncryptedSecret string
	UserID          int64
	Name            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	URI             string
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
	results := make([]*ClientResult, 0)
	err := rq.Select(cq.ID, cq.EncryptedSecret, cq.UserID, cq.Name, cq.CreatedAt, cq.UpdatedAt, rq.URI).
		Where(c.query.Client.ID.Eq(string(id))).LeftJoin(cq, rq.ClientID.EqCol(cq.ID)).Scan(&results)
	if err != nil {
		return nil, err
	}
	return toDomainClient(results), nil
}

func toDomainClient(results []*ClientResult) *oauth.Client {
	uris := make([]string, 0, len(results))
	for _, r := range results {
		uris = append(uris, r.URI)
	}
	return &oauth.Client{
		ID:              oauth.ClientID(results[0].ID),
		EncryptedSecret: results[0].EncryptedSecret,
		UserID:          results[0].UserID,
		Name:            results[0].Name,
		CreatedAt:       results[0].CreatedAt,
		UpdatedAt:       results[0].UpdatedAt,
		RedirectURIs:    uris,
	}
}
