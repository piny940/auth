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
	Name              string
	EncryptedPassword string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
type IUserRepo interface {
	FindByID(id int64) (*User, error)
	FindByName(name string) (*User, error)
	Create(name, encryptedPassword string) error
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
	hash, err := EncryptPassword(password)
	if err != nil {
		return false, err
	}
	return u.EncryptedPassword == hash, nil
}

func (s *UserService) Validate(name, password, passwordConfirmation string) error {
	if len(name) < MIN_NAME_LENGTH {
		return ErrNameLengthNotEnough{}
	}
	_, err := s.UserRepo.FindByName(name)
	if err == nil {
		return ErrNameAlreadyUsed{}
	} else if !errors.Is(err, ErrRecordNotFound{}) {
		return err
	}
	if !passwordStrong(password) {
		return ErrPasswordLengthNotEnough{}
	}
	if password != passwordConfirmation {
		return ErrPasswordConfirmation{}
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

type ErrPasswordConfirmation struct{}

func (e ErrPasswordConfirmation) Error() string {
	return "password and password confirmation do not match"
}

type ErrNameAlreadyUsed struct{}

func (e ErrNameAlreadyUsed) Error() string {
	return "this name is already used"
}

type ErrNameLengthNotEnough struct{}

func (e ErrNameLengthNotEnough) Error() string {
	return fmt.Sprintf("name length must be at least %d", MIN_NAME_LENGTH)
}

type ErrPasswordLengthNotEnough struct{}

func (e ErrPasswordLengthNotEnough) Error() string {
	return fmt.Sprintf("password length must be at least %d", MIN_PASSWORD_LENGTH)
}
