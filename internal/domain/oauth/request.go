package oauth

import (
	"auth/internal/domain"
	"crypto/rand"
	"errors"
	"slices"
	"time"
)

type TypeResponseType string

const (
	ResponseTypeCode TypeResponseType = "code"
)

var AllResponseTypes = []TypeResponseType{
	ResponseTypeCode,
}

type TypeScope string

const (
	ScopeOpenID TypeScope = "openid"
)

var AllScopes = []TypeScope{
	ScopeOpenID,
}

type AuthRequest struct {
	ResponseType TypeResponseType
	ClientID     ClientID
	RedirectURI  string
	Scopes       []TypeScope
	State        *string
}

type AuthService struct {
	ClientRepo   IClientRepo
	ApprovalRepo IApprovalRepo
	AuthCodeRepo IAuthCodeRepo
}

const (
	AUTH_CODE_TTL = 5 * time.Minute
	AUTH_CODE_LEN = 32
)

func NewAuthService(clientRepo IClientRepo, approvalRepo IApprovalRepo, authCodeRepo IAuthCodeRepo) *AuthService {
	return &AuthService{
		ClientRepo:   clientRepo,
		ApprovalRepo: approvalRepo,
		AuthCodeRepo: authCodeRepo,
	}
}

func (s *AuthService) Validate(r *AuthRequest) error {
	if !slices.Contains(AllResponseTypes, r.ResponseType) {
		return ErrInvalidRequestType
	}
	client, err := s.ClientRepo.FindByID(r.ClientID)
	if errors.Is(err, domain.ErrRecordNotFound) {
		return ErrInvalidClientID
	}
	if err != nil {
		return err
	}
	if !slices.Contains(client.RedirectURIs, r.RedirectURI) {
		return ErrInvalidRedirectURI
	}
	if err := ValidScopes(r.Scopes); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Approved(r *AuthRequest, user *domain.User) (bool, error) {
	approval, err := s.ApprovalRepo.Find(r.ClientID, user.ID)
	if errors.Is(err, domain.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if approval == nil {
		return false, nil
	}
	for _, scope := range r.Scopes {
		if !slices.Contains(approval.Scopes, scope) {
			return false, nil
		}
	}
	return true, nil
}

func (s *AuthService) IssueAuthCode(clientID ClientID, userID domain.UserID, scopes []TypeScope) (*AuthCode, error) {
	v := make([]byte, AUTH_CODE_LEN)
	if _, err := rand.Read(v); err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(AUTH_CODE_TTL)
	if err := s.AuthCodeRepo.Create(string(v), clientID, userID, scopes, expiresAt); err != nil {
		return nil, err
	}
	return &AuthCode{
		Value:     string(v),
		ClientID:  clientID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Scopes:    scopes,
	}, nil
}

func ValidScopes(scopes []TypeScope) error {
	for _, s := range scopes {
		if !slices.Contains(AllScopes, s) {
			return ErrInvalidScope
		}
	}
	return nil
}

var (
	ErrInvalidRequestType = errors.New("invalid request type")
	ErrInvalidClientID    = errors.New("invalid client id. client not found")
	ErrInvalidRedirectURI = errors.New("invalid redirect uri")
	ErrInvalidScope       = errors.New("invalid scope")
)
