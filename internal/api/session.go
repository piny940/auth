package api

import (
	"auth/internal/domain"
	"net/http"

	"github.com/gorilla/sessions"
)

const SESSION_NAME = "auth_piny940"

var sessionsOptions = &sessions.Options{
	HttpOnly: true,
	Secure:   true,
	MaxAge:   60 * 60 * 24 * 7,
}

var store sessions.Store

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

func Login(r *http.Request, w http.ResponseWriter, user *domain.User) error {
	return setToSession(r, w, SESSION_USER_KEY, user)
}
func Logout(r *http.Request, w http.ResponseWriter) error {
	return setToSession(r, w, SESSION_USER_KEY, nil)
}
func CurrentUser(r *http.Request) (*domain.User, error) {
	user, err := getFromSession(r, SESSION_USER_KEY)
	if err != nil {
		return nil, err
	}
	u, ok := user.(*domain.User)
	if !ok {
		return nil, nil
	}
	return u, nil
}
