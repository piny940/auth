package usecase

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
)

type AuthUsecase struct {
	UserRepo     domain.IUserRepo
	ApprovalRepo oauth.IApprovalRepo
	UserService  *domain.UserService
}

func NewAuthUsecase(
	userRepo domain.IUserRepo,
	approvalRepo oauth.IApprovalRepo,
	userSvc *domain.UserService,
) *AuthUsecase {
	return &AuthUsecase{
		UserRepo:     userRepo,
		ApprovalRepo: approvalRepo,
		UserService:  userSvc,
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
		return nil, ErrPasswordNotMatch{}
	}
	return user, nil
}

type ErrPasswordNotMatch struct{}

func (e ErrPasswordNotMatch) Error() string {
	return "invalid password"
}

func (u *AuthUsecase) SignUp(username, password, passwordConfirmation string) error {
	if err := u.UserService.Validate(username, password, passwordConfirmation); err != nil {
		return err
	}
	hash, err := domain.EncryptPassword(password)
	if err != nil {
		return err
	}
	err = u.UserRepo.Create(username, hash)
	if err != nil {
		return err
	}
	return nil
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
