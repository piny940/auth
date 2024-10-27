package oauth

import (
	"auth/internal/domain"
	"crypto/rand"
	"errors"
	"time"
)

type AuthCode struct {
	Value       string
	ClientID    ClientID
	UserID      domain.UserID
	ExpiresAt   time.Time
	Used        bool
	RedirectURI string
	Scopes      []TypeScope
}
type IAuthCodeRepo interface {
	Find(value string) (*AuthCode, error)
	Create(value string, clientID ClientID, userID domain.UserID, scopes []TypeScope, expiresAt time.Time, redirectURI string) error
}

type AuthCodeService struct {
	AuthCodeRepo IAuthCodeRepo
}

func NewAuthCodeService(authCodeRepo IAuthCodeRepo) *AuthCodeService {
	return &AuthCodeService{
		AuthCodeRepo: authCodeRepo,
	}
}

func (s *AuthCodeService) IssueAuthCode(clientID ClientID, userID domain.UserID, scopes []TypeScope, redirectURI string) (*AuthCode, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, AUTH_CODE_LEN)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	var code string
	for _, v := range b {
		code += string(letters[int(v)%len(letters)])
	}
	expiresAt := time.Now().Add(AUTH_CODE_TTL)
	if err := s.AuthCodeRepo.Create(code, clientID, userID, scopes, expiresAt, redirectURI); err != nil {
		return nil, err
	}
	return &AuthCode{
		Value:     code,
		ClientID:  clientID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Scopes:    scopes,
	}, nil
}

func (s *AuthCodeService) Verify(code string, clientID ClientID, redirectURI string) (*AuthCode, error) {
	authCode, err := s.AuthCodeRepo.Find(code)
	if errors.Is(err, domain.ErrRecordNotFound) {
		return nil, domain.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	if authCode.ClientID != clientID {
		return nil, ErrInvalidClientID
	}
	if authCode.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredAuthCode
	}
	if authCode.Used {
		return nil, ErrUsedAuthCode
	}
	if authCode.RedirectURI != redirectURI {
		return nil, ErrInvalidRedirectURI
	}
	return authCode, nil
}

var (
	ErrExpiredAuthCode = errors.New("expired auth code")
	ErrUsedAuthCode    = errors.New("used auth code")
)
