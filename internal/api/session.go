package api

import (
	"auth/internal/domain"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const SESSION_NAME = "com.piny940.auth"

type sessionStore struct {
	store sessions.Store
}

var SessionStore *sessionStore = &sessionStore{
	sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET"))),
}
var sessionsOptions = &sessions.Options{
	HttpOnly: true,
	Secure:   true,
	MaxAge:   60 * 60 * 24 * 7,
}

func (s *sessionStore) Get(r *http.Request, key string) (interface{}, error) {
	session, err := s.store.Get(r, SESSION_NAME)
	if err != nil {
		return nil, err
	}
	return session.Values[key], nil
}

func (s *sessionStore) Set(r *http.Request, w http.ResponseWriter, key string, value interface{}) error {
	session, err := s.store.Get(r, SESSION_NAME)
	if err != nil {
		return err
	}
	session.Options = sessionsOptions
	session.Values[key] = value
	return session.Save(r, w)
}

const SESSION_USER_KEY = "user"

func Login(r *http.Request, w http.ResponseWriter, user *domain.User) error {
	return SessionStore.Set(r, w, SESSION_USER_KEY, user)
}
func CurrentUser(r *http.Request) (*domain.User, error) {
	user, err := SessionStore.Get(r, SESSION_USER_KEY)
	if err != nil {
		return nil, err
	}
	u, ok := user.(*domain.User)
	if !ok {
		return nil, nil
	}
	return u, nil
}
