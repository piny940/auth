package usecase

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"errors"
	"fmt"
	"slices"

	"github.com/lestrrat-go/jwx/jwk"
)

type OAuthUsecase struct {
	AuthCodeService *oauth.AuthCodeService
	ApprovalService *oauth.ApprovalService
	TokenService    *oauth.TokenService
	ApprovalRepo    oauth.IApprovalRepo
	JWKsService     *oauth.JWKsService
	ClientRepo      oauth.IClientRepo
}

func NewOAuthUsecase(
	authCodeSvc *oauth.AuthCodeService,
	jwksService *oauth.JWKsService,
	approvalSvc *oauth.ApprovalService,
	approvalRepo oauth.IApprovalRepo,
	tokenSvc *oauth.TokenService,
	clientRepo oauth.IClientRepo,
) *OAuthUsecase {
	return &OAuthUsecase{
		JWKsService:     jwksService,
		AuthCodeService: authCodeSvc,
		ApprovalService: approvalSvc,
		ApprovalRepo:    approvalRepo,
		TokenService:    tokenSvc,
		ClientRepo:      clientRepo,
	}
}

type TypeResponseType string

const (
	ResponseTypeCode TypeResponseType = "code"
)

var AllResponseTypes = []TypeResponseType{
	ResponseTypeCode,
}

type AuthRequest struct {
	ResponseType TypeResponseType
	ClientID     oauth.ClientID
	RedirectURI  string
	Scopes       []oauth.TypeScope
	State        *string
}

func (u *OAuthUsecase) RequestCodeAuth(user *domain.User, req *AuthRequest) (*oauth.AuthCode, error) {
	if req.ResponseType != ResponseTypeCode {
		return nil, ErrInvalidRequestType
	}
	client, err := u.ClientRepo.FindByID(req.ClientID)
	if errors.Is(err, domain.ErrRecordNotFound) {
		return nil, ErrClientNotFound
	}
	if !client.RedirectURIValid(req.RedirectURI) {
		return nil, ErrRedirectURINotRegistered
	}
	if err := oauth.ValidScopes(req.Scopes); err != nil {
		return nil, fmt.Errorf("invalid scopes: %w", err)
	}
	ok, err := u.ApprovalService.Approved(req.ClientID, user.ID, req.Scopes)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotApproved
	}
	return u.AuthCodeService.IssueAuthCode(req.ClientID, user.ID, req.Scopes, req.RedirectURI)
}

func (u *OAuthUsecase) Approve(user *domain.User, clientID oauth.ClientID, scopes []oauth.TypeScope) error {
	return u.ApprovalRepo.Approve(clientID, user.ID, scopes)
}

type TokenRequest struct {
	GrantType   string
	AuthCode    string
	RedirectURI string
	ClientID    oauth.ClientID
}
type TypeGrantType string

const (
	GrantTypeAuthorizationCode TypeGrantType = "authorization_code"
)

func (u *OAuthUsecase) RequestToken(req *TokenRequest) (*oauth.AccessToken, *oauth.IDToken, error) {
	if req.GrantType != string(GrantTypeAuthorizationCode) {
		return nil, nil, ErrInvalidGrantType
	}
	authCode, err := u.AuthCodeService.Verify(req.AuthCode, req.ClientID, req.RedirectURI)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid auth code: %w", err)
	}
	if err := u.AuthCodeService.AuthCodeRepo.Use(authCode.Value); err != nil {
		return nil, nil, fmt.Errorf("failed to use auth code: %w", err)
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

func (u *OAuthUsecase) GetJWKs() (jwk.Set, error) {
	set, err := u.JWKsService.IssueJwks()
	if err != nil {
		return nil, fmt.Errorf("failed to issue jwks: %w", err)
	}
	return set, nil
}

var (
	ErrInvalidRequestType       = errors.New("invalid request type")
	ErrClientNotFound           = errors.New("invalid client id. client not found")
	ErrRedirectURINotRegistered = errors.New("invalid redirect uri. this redirect uri is not registered")
	ErrPasswordNotMatch         = errors.New("invalid password")
	ErrNotApproved              = errors.New("not approved")
	ErrInvalidGrantType         = errors.New("invalid grant type")
)
