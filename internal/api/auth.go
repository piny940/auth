package api

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"context"
	"net/http"
	"time"
)

type AuthSession struct {
	User     *domain.User
	AuthTime time.Time
}

type Auth interface {
	CurrentUser(ctx context.Context) (*AuthSession, error)
	Login(ctx context.Context, user *domain.User) (*http.Cookie, error)
	Logout(ctx context.Context) (*http.Cookie, error)
	AccessScopes(ctx context.Context) ([]oauth.TypeScope, error)
}
