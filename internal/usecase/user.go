package usecase

import (
	"auth/internal/domain"
)

type UserUsecase struct {
	UserService *domain.UserService
	UserRepo    domain.IUserRepo
}

func NewAuthUsecase(
	userSvc *domain.UserService,
	userRepo domain.IUserRepo,
) *UserUsecase {
	return &UserUsecase{
		UserService: userSvc,
		UserRepo:    userRepo,
	}
}

func (u *UserUsecase) Login(username, password string) (*domain.User, error) {
	user, err := u.UserRepo.FindByName(username)
	if err != nil {
		return nil, err
	}
	ok, err := user.PasswordMatch(password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrPasswordNotMatch
	}
	return user, nil
}

func (u *UserUsecase) SignUp(email, username, password, passwordConfirmation string) (*domain.User, error) {
	if err := u.UserService.Validate(username, password, passwordConfirmation); err != nil {
		return nil, err
	}
	hash, err := domain.EncryptPassword(password)
	if err != nil {
		return nil, err
	}
	err = u.UserRepo.Create(email, username, hash)
	if err != nil {
		return nil, err
	}
	user, err := u.UserRepo.FindByName(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
