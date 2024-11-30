package api

import (
	"auth/internal/domain"
	"encoding/gob"
)

func Init() {
	gob.Register(&domain.User{})
}
