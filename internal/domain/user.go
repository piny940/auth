package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int64
	Name              string
	EncryptedPassword string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
type IUserRepo interface {
	FindByID(id int64) (*User, error)
	FindByName(name string) (*User, error)
}

func (u *User) ValidPassword(password string) (bool, error) {
	hash, err := u.encryptPassword(password)
	if err != nil {
		return false, err
	}
	return u.EncryptedPassword == hash, nil
}

func (u *User) encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
