package oauth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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

func (c *Client) SecretCorrect(secret string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(c.EncryptedSecret), []byte(secret))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func EncryptClientSecret(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
