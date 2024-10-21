package oauth

import "time"

type ClientID string
type Client struct {
	ID              ClientID
	EncryptedSecret string
	UserID          int64
	Name            string
	RedirectURIs    []string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
type IClientRepo interface {
	FindByID(id ClientID) (*Client, error)
}
