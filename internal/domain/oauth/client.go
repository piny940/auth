package oauth

import (
	"auth/internal/domain"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ClientID string
type Client struct {
	ID              ClientID
	EncryptedSecret string
	UserID          domain.UserID
	Name            string
	RedirectURIs    []string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
type ClientInput struct {
	ID              ClientID
	Name            string
	EncryptedSecret string
	UserID          domain.UserID
	RedirectURIs    []string
}
type IClientRepo interface {
	FindByID(id ClientID) (*Client, error)
	FindWithUserID(id ClientID, userID domain.UserID) (*Client, error)
	List(userID domain.UserID) ([]*Client, error)
	Create(client *ClientInput) error
	Update(client *Client, userID domain.UserID) error
	Delete(id ClientID, userID domain.UserID) error
}

const CLIENT_ID_LEN = 16

func IssueClientID() (ClientID, error) {
	str, err := randomString(CLIENT_ID_LEN)
	if err != nil {
		return "", err
	}
	return ClientID(str), nil
}
func IssueClientSecret() (string, string, error) {
	raw, err := randomString(32)
	if err != nil {
		return "", "", err
	}
	encrypted, err := EncryptClientSecret(raw)
	if err != nil {
		return "", "", err
	}
	return raw, encrypted, nil
}

func (c *Client) SecretCorrect(secret string) error {
	err := bcrypt.CompareHashAndPassword([]byte(c.EncryptedSecret), []byte(secret))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrInvalidClientSecret
	}
	if err != nil {
		return err
	}
	return nil
}

func EncryptClientSecret(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

var (
	ErrInvalidClientSecret = errors.New("invalid client secret")
)
