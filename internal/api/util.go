package api

import (
	"auth/internal/domain"
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const SESSION_NAME = "auth"

var sessionsOptions = &sessions.Options{
	HttpOnly: true,
	Secure:   true,
	MaxAge:   60 * 60 * 24 * 7,
}

var store *sessions.CookieStore

type CtxKey string

const (
	CtxKeyUser CtxKey = "user"
)

func getFromSession(r *http.Request, key string) (interface{}, error) {
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		return nil, err
	}
	return session.Values[key], nil
}

func setToSession(r *http.Request, w http.ResponseWriter, key string, value interface{}) error {
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		return err
	}
	session.Options = sessionsOptions
	session.Values[key] = value
	return session.Save(r, w)
}

const SESSION_USER_KEY = "user"

func Login(c context.Context, user *domain.User) (*http.Cookie, error) {
	encoded, err := securecookie.EncodeMulti(SESSION_NAME, map[interface{}]interface{}{
		SESSION_USER_KEY: user,
	}, store.Codecs...)
	if err != nil {
		return nil, err
	}
	return sessions.NewCookie(SESSION_NAME, encoded, sessionsOptions), nil
}
func Logout(r *http.Request, w http.ResponseWriter) error {
	return setToSession(r, w, SESSION_USER_KEY, nil)
}

func SetCurrentUserToCtx(c echo.Context) error {
	user, err := UserFromReq(c.Request())
	if err != nil {
		return err
	}
	ctx := context.WithValue(c.Request().Context(), CtxKeyUser, user)
	req := c.Request().WithContext(ctx)
	c.SetRequest(req)
	return nil
}
func CurrentUser(c context.Context) (*domain.User, error) {
	user, ok := c.Value(CtxKeyUser).(*domain.User)
	if !ok {
		return nil, ErrUnauthorized
	}
	return user, nil
}
func UserFromReq(r *http.Request) (*domain.User, error) {
	user, err := getFromSession(r, SESSION_USER_KEY)
	if err != nil {
		return nil, err
	}
	u, ok := user.(*domain.User)
	if !ok {
		return nil, ErrUnauthorized
	}
	return u, nil
}

var ErrUnauthorized = errors.New("unauthorized")
