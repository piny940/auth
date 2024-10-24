package oauth

import (
	"auth/internal/domain"
	"crypto/rand"
	"time"
)

type AuthCode struct {
	Value     string
	ClientID  ClientID
	UserID    domain.UserID
	ExpiresAt time.Time
	Scopes    []TypeScope
}
type IAuthCodeRepo interface {
	Create(value string, clientID ClientID, userID domain.UserID, scopes []TypeScope, expiresAt time.Time) error
}

type AuthCodeService struct {
	AuthCodeRepo IAuthCodeRepo
}

func (s *AuthCodeService) IssueAuthCode(clientID ClientID, userID domain.UserID, scopes []TypeScope) (*AuthCode, error) {
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
	if err := s.AuthCodeRepo.Create(code, clientID, userID, scopes, expiresAt); err != nil {
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
