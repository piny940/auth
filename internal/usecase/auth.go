package usecase

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"errors"
)

type AuthUsecase struct {
	UserRepo     domain.IUserRepo
	ApprovalRepo oauth.IApprovalRepo
	UserService  *domain.UserService
	AuthService  *oauth.AuthService
}

func NewAuthUsecase(
	userRepo domain.IUserRepo,
	approvalRepo oauth.IApprovalRepo,
	userSvc *domain.UserService,
	authSvc *oauth.AuthService,
) *AuthUsecase {
	return &AuthUsecase{
		UserRepo:     userRepo,
		ApprovalRepo: approvalRepo,
		UserService:  userSvc,
		AuthService:  authSvc,
	}
}

func (u *AuthUsecase) Login(username, password string) (*domain.User, error) {
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

func (u *AuthUsecase) SignUp(username, password, passwordConfirmation string) (*domain.User, error) {
	if err := u.UserService.Validate(username, password, passwordConfirmation); err != nil {
		return nil, err
	}
	hash, err := domain.EncryptPassword(password)
	if err != nil {
		return nil, err
	}
	err = u.UserRepo.Create(username, hash)
	if err != nil {
		return nil, err
	}
	user, err := u.UserRepo.FindByName(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *AuthUsecase) Request(user *domain.User, req *oauth.AuthRequest) error {
	err := u.AuthService.Validate(req)
	if err != nil {
		return err
	}
	ok, err := u.AuthService.Approved(req, user)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotApproved
	}
	return nil
}

var (
	ErrPasswordNotMatch = errors.New("invalid password")
	ErrNotApproved      = errors.New("not approved")
)
