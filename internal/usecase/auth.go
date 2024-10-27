package usecase

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"errors"
	"fmt"
)

type AuthUsecase struct {
	UserService     *domain.UserService
	RequestService  *oauth.RequestService
	AuthCodeService *oauth.AuthCodeService
	ApprovalService *oauth.ApprovalService
	TokenService    *oauth.TokenService
	UserRepo        domain.IUserRepo
	ApprovalRepo    oauth.IApprovalRepo
	AuthCodeRepo    oauth.IAuthCodeRepo
	ClientRepo      oauth.IClientRepo
}

func NewAuthUsecase(
	userSvc *domain.UserService,
	requestSvc *oauth.RequestService,
	authCodeSvc *oauth.AuthCodeService,
	approvalSvc *oauth.ApprovalService,
	tokenSvc *oauth.TokenService,
	userRepo domain.IUserRepo,
	approvalRepo oauth.IApprovalRepo,
	authCodeRepo oauth.IAuthCodeRepo,
	clientRepo oauth.IClientRepo,
) *AuthUsecase {
	return &AuthUsecase{
		UserService:     userSvc,
		RequestService:  requestSvc,
		AuthCodeService: authCodeSvc,
		ApprovalService: approvalSvc,
		TokenService:    tokenSvc,
		UserRepo:        userRepo,
		ApprovalRepo:    approvalRepo,
		AuthCodeRepo:    authCodeRepo,
		ClientRepo:      clientRepo,
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
	return u.AuthCodeService.IssueAuthCode(req.ClientID, user.ID, req.Scopes, req.RedirectURI)
}

func (u *AuthUsecase) Approve(user *domain.User, clientID oauth.ClientID, scopes []oauth.TypeScope) error {
	return u.ApprovalService.Approve(clientID, user.ID, scopes)
}

type TokenRequest struct {
	GrantType    string
	AuthCode     string
	RedirectURI  string
	ClientID     oauth.ClientID
	ClientSecret string
}
type TypeGrantType string

const (
	GrantTypeAuthorizationCode TypeGrantType = "authorization_code"
)

func (u *AuthUsecase) GetToken(req *TokenRequest) (*oauth.AccessToken, error) {
	if req.GrantType != string(GrantTypeAuthorizationCode) {
		return nil, ErrInvalidGrantType
	}
	client, err := u.ClientRepo.FindByID(req.ClientID)
	if err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}
	if err := client.SecretCorrect(req.ClientSecret); err != nil {
		return nil, fmt.Errorf("invalid client secret: %w", err)
	}
	authCode, err := u.AuthCodeService.Verify(req.AuthCode, req.ClientID, req.RedirectURI)
	if err != nil {
		return nil, fmt.Errorf("invalid auth code: %w", err)
	}
	accessToken, err := u.TokenService.IssueAccessToken(authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to issue access token: %w", err)
	}
	return accessToken, nil
}

var (
	ErrPasswordNotMatch = errors.New("invalid password")
	ErrNotApproved      = errors.New("not approved")
	ErrInvalidGrantType = errors.New("invalid grant type")
)
