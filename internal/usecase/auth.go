package usecase

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
)

type AuthUsecase struct {
	UserRepo     domain.IUserRepo
	ApprovalRepo oauth.IApprovalRepo
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

func (u *AuthUsecase) Request(user *domain.User, req *oauth.AuthRequest) error {
	err := req.Validate()
	if err != nil {
		return err
	}
	ok, err := req.ApprovedBy(user)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotApproved{}
	}
	return nil
}

type ErrNotApproved struct{}

func (e ErrNotApproved) Error() string {
	return "not approved"
}
