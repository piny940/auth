package oauth

import "time"

type Client struct {
	ID              string
	EncryptedSecret string
	UserID          int64
	Name            string
	RedirectURIs    []string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
type ClientRepo interface {
	FindByID(id string) (*Client, error)
}
