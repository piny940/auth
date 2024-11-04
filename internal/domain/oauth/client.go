package oauth

import (
	"auth/internal/domain"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
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
const CLIENT_SECRET_LEN = 32

func IssueClientID() (ClientID, error) {
	str, err := randomString(CLIENT_ID_LEN)
	if err != nil {
		return "", err
	}
	return ClientID(str), nil
}
func IssueClientSecret() (string, string, error) {
	raw, err := randomString(CLIENT_SECRET_LEN)
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

func (c *Client) RedirectURIValid(rawURL string) bool {
	url, err := url.Parse(rawURL)
	if err != nil {
		slog.Info(fmt.Sprintf("failed to parse uri: %s", rawURL))
		return false
	}
	if url.Fragment != "" {
		slog.Info(fmt.Sprintf("fragment is not allowed: %s", rawURL))
		return false
	}
	for _, clientUri := range c.RedirectURIs {
		clientURL, err := url.Parse(clientUri)
		if err != nil {
			slog.Info(fmt.Sprintf("failed to parse uri: %s", clientUri))
			continue
		}
		if url.Scheme != clientURL.Scheme {
			continue
		}
		if url.Host != clientURL.Host {
			continue
		}
		if url.Path != clientURL.Path {
			continue
		}
		return true
	}
	return false
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
