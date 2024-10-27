package usecase

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"errors"
	"fmt"
	"slices"
)

type OAuthUsecase struct {
	RequestService  *oauth.RequestService
	AuthCodeService *oauth.AuthCodeService
	ApprovalService *oauth.ApprovalService
	TokenService    *oauth.TokenService
	ClientRepo      oauth.IClientRepo
}

func NewOAuthUsecase(
	reqSvc *oauth.RequestService,
	authCodeSvc *oauth.AuthCodeService,
	approvalSvc *oauth.ApprovalService,
	tokenSvc *oauth.TokenService,
	clientRepo oauth.IClientRepo,
) *OAuthUsecase {
	return &OAuthUsecase{
		RequestService:  reqSvc,
		AuthCodeService: authCodeSvc,
		ApprovalService: approvalSvc,
		TokenService:    tokenSvc,
		ClientRepo:      clientRepo,
	}
}

func (u *OAuthUsecase) RequestAuthorization(user *domain.User, req *oauth.AuthRequest) (*oauth.AuthCode, error) {
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

func (u *OAuthUsecase) Approve(user *domain.User, clientID oauth.ClientID, scopes []oauth.TypeScope) error {
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

func (u *OAuthUsecase) RequestToken(req *TokenRequest) (*oauth.AccessToken, *oauth.IDToken, error) {
	if req.GrantType != string(GrantTypeAuthorizationCode) {
		return nil, nil, ErrInvalidGrantType
	}
	client, err := u.ClientRepo.FindByID(req.ClientID)
	if err != nil {
		return nil, nil, fmt.Errorf("client not found: %w", err)
	}
	if err := client.SecretCorrect(req.ClientSecret); err != nil {
		return nil, nil, fmt.Errorf("invalid client secret: %w", err)
	}
	authCode, err := u.AuthCodeService.Verify(req.AuthCode, req.ClientID, req.RedirectURI)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid auth code: %w", err)
	}
	accessToken, err := u.TokenService.IssueAccessToken(authCode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to issue access token: %w", err)
	}
	if slices.Contains(authCode.Scopes, oauth.ScopeOpenID) {
		idToken, err := u.TokenService.IssueIDToken(authCode)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to issue id token: %w", err)
		}
		return accessToken, idToken, nil
	} else {
		return accessToken, nil, nil
	}
}

var (
	ErrPasswordNotMatch = errors.New("invalid password")
	ErrNotApproved      = errors.New("not approved")
	ErrInvalidGrantType = errors.New("invalid grant type")
)
