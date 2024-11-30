package api

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthSession struct {
	User     *domain.User
	AuthTime time.Time
}

type Auth interface {
	CurrentUser(ctx context.Context) (*AuthSession, error)
	Login(ctx context.Context, user *domain.User) error
	Logout(ctx context.Context) error
	AccessScopes(ctx context.Context) ([]oauth.TypeScope, error)
}

type EchoContextReg interface {
	Context(c context.Context) (echo.Context, error)
}

func toQueryString(h map[string]string) string {
	var q []string
	for k, v := range h {
		q = append(q, k+"="+url.QueryEscape(v))
	}
	return strings.Join(q, "&")
}

var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrNotFoundInSession = errors.New("not found in session")
)

func ptr[T any](v T) *T {
	return &v
}
