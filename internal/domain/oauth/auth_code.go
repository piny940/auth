package oauth

import (
	"auth/internal/domain"
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
