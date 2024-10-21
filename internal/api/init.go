package api

import (
	"auth/internal/domain"
	"encoding/gob"
	"os"

	"github.com/gorilla/sessions"
)

func Init() {
	gob.Register(&domain.User{})
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
}
