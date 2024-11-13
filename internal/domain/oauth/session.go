package oauth

import (
	"auth/internal/domain"
	"errors"
	"time"
)

type Session struct {
	User     *domain.User
	AuthTime time.Time
}

func (s *Session) Valid(maxAge time.Duration) error {
	if time.Since(s.AuthTime) > maxAge {
		return ErrSessionExpired
	}
	return nil
}

var (
	ErrSessionExpired = errors.New("session expired")
)
