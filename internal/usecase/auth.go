package usecase

import "auth/internal/domain"

type AuthUsecase struct {
	UserRepo domain.IUserRepo
}

func NewAuthUsecase(userRepo domain.IUserRepo) *AuthUsecase {
	return &AuthUsecase{
		UserRepo: userRepo,
	}
}

func (u *AuthUsecase) Login(username, password string) (*domain.User, error) {
	user, err := u.UserRepo.FindByName(username)
	if err != nil {
		return nil, err
	}
	ok, err := user.ValidPassword(password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidPassword{}
	}
	return user, nil
}

type ErrInvalidPassword struct{}

func (e ErrInvalidPassword) Error() string {
	return "invalid password"
}
