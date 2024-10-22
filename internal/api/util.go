package api

import (
	"auth/internal/domain"
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

const SESSION_NAME = "auth"

var sessionsOptions = &sessions.Options{
	HttpOnly: true,
	SameSite: http.SameSiteLaxMode,
	// Secure:   true,
	MaxAge: 60 * 60 * 24 * 7,
}

var store *sessions.CookieStore

type CtxKey string

const (
	CtxKeySessionRegistry CtxKey = "sessionRegistry"
)

const SESSION_USER_KEY = "user"

func SetSessionRegistry(r *http.Request) error {
	reg := sessions.GetRegistry(r)
	ctx := context.WithValue(r.Context(), CtxKeySessionRegistry, reg)
	*r = *r.WithContext(ctx)
	return nil
}

func GetFromSession(c context.Context, key string) (interface{}, error) {
	reg := c.Value(CtxKeySessionRegistry).(*sessions.Registry)
	session, err := reg.Get(store, SESSION_NAME)
	if err != nil {
		return nil, err
	}
	v, ok := session.Values[key]
	if !ok {
		return nil, ErrNotFoundInSession
	}
	return v, nil
}
func SetToSession(c context.Context, key string, value interface{}) error {
	reg := c.Value(CtxKeySessionRegistry).(*sessions.Registry)
	session, err := reg.Get(store, SESSION_NAME)
	if err != nil {
		return err
	}
	session.Options = sessionsOptions
	session.Values[key] = value
	return nil
}
func Save(c context.Context) (*http.Cookie, error) {
	reg := c.Value(CtxKeySessionRegistry).(*sessions.Registry)
	session, err := reg.Get(store, SESSION_NAME)
	if err != nil {
		return nil, err
	}
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values, store.Codecs...)
	if err != nil {
		return nil, err
	}
	return sessions.NewCookie(session.Name(), encoded, sessionsOptions), nil
}

func CurrentUser(c context.Context) (*domain.User, error) {
	userObj, err := GetFromSession(c, SESSION_USER_KEY)
	if errors.Is(err, ErrNotFoundInSession) {
		return nil, ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}
	user, ok := userObj.(*domain.User)
	if !ok {
		return nil, ErrUnauthorized
	}
	return user, nil
}

func Login(c context.Context, user *domain.User) (*http.Cookie, error) {
	err := SetToSession(c, SESSION_USER_KEY, user)
	if err != nil {
		return nil, err
	}
	return Save(c)
}

func Logout(c context.Context) (*http.Cookie, error) {
	err := SetToSession(c, SESSION_USER_KEY, nil)
	if err != nil {
		return nil, err
	}
	return Save(c)
}

var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrNotFoundInSession = errors.New("not found in session")
)
