package domain

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserID int64
type User struct {
	ID                UserID
	Email             string
	Name              string
	EncryptedPassword string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
type IUserRepo interface {
	FindByID(id UserID) (*User, error)
	FindByName(name string) (*User, error)
	Create(email, name, encryptedPassword string) error
}

type UserService struct {
	UserRepo IUserRepo
}

const (
	MIN_NAME_LENGTH     = 3
	MIN_PASSWORD_LENGTH = 8
)

func NewUserService(userRepo IUserRepo) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (u *User) PasswordMatch(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *UserService) Validate(name, password, passwordConfirmation string) error {
	if len(name) < MIN_NAME_LENGTH {
		return ErrNameLengthNotEnough
	}
	_, err := s.UserRepo.FindByName(name)
	if err == nil {
		return ErrNameAlreadyUsed
	} else if !errors.Is(err, ErrRecordNotFound) {
		return err
	}
	if !passwordStrong(password) {
		return ErrPasswordLengthNotEnough
	}
	if password != passwordConfirmation {
		return ErrPasswordConfirmation
	}
	return nil
}

func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func passwordStrong(password string) bool {
	return len(password) >= MIN_PASSWORD_LENGTH
}

var (
	ErrPasswordConfirmation    = errors.New("password and password confirmation do not match")
	ErrNameAlreadyUsed         = errors.New("this name is already used")
	ErrNameLengthNotEnough     = fmt.Errorf("name length must be at least %d", MIN_NAME_LENGTH)
	ErrPasswordLengthNotEnough = fmt.Errorf("password length must be at least %d", MIN_PASSWORD_LENGTH)
)
