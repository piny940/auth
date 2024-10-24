package usecase

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"errors"
)

type AuthUsecase struct {
	UserService     *domain.UserService
	RequestService  *oauth.RequestService
	AuthCodeService *oauth.AuthCodeService
	ApprovalService *oauth.ApprovalService
	UserRepo        domain.IUserRepo
	ApprovalRepo    oauth.IApprovalRepo
	AuthCodeRepo    oauth.IAuthCodeRepo
}

func NewAuthUsecase(
	userSvc *domain.UserService,
	requestSvc *oauth.RequestService,
	authCodeSvc *oauth.AuthCodeService,
	approvalSvc *oauth.ApprovalService,
	userRepo domain.IUserRepo,
	approvalRepo oauth.IApprovalRepo,
	authCodeRepo oauth.IAuthCodeRepo,
) *AuthUsecase {
	return &AuthUsecase{
		UserService:     userSvc,
		RequestService:  requestSvc,
		AuthCodeService: authCodeSvc,
		ApprovalService: approvalSvc,
		UserRepo:        userRepo,
		ApprovalRepo:    approvalRepo,
		AuthCodeRepo:    authCodeRepo,
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

func (u *AuthUsecase) Request(user *domain.User, req *oauth.AuthRequest) (*oauth.AuthCode, error) {
	err := u.RequestService.Validate(req)
	if err != nil {
		return nil, err
	}
	ok, err := u.ApprovalService.Approved(req, user)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotApproved
	}
	return u.AuthCodeService.IssueAuthCode(req.ClientID, user.ID, req.Scopes)
}

func (u *AuthUsecase) Approve(user *domain.User, clientID oauth.ClientID, scopes []oauth.TypeScope) error {
	_, err := u.AuthService.ClientRepo.FindByID(clientID)
	if err != nil {
		return err
	}
	if err := oauth.ValidScopes(scopes); err != nil {
		return err
	}
	if err := u.AuthService.ApprovalRepo.Create(clientID, user.ID, scopes); err != nil {
		return err
	}
	return nil
}

var (
	ErrPasswordNotMatch = errors.New("invalid password")
	ErrNotApproved      = errors.New("not approved")
)
