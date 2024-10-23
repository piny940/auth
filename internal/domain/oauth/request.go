package oauth

import (
	"auth/internal/domain"
	"errors"
	"slices"
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

type IApprovalRepo interface {
	Find(clientID ClientID, userID domain.UserID) (*Approval, error)
}

type AuthService struct {
	ClientRepo   IClientRepo
	ApprovalRepo IApprovalRepo
}

func NewAuthService(clientRepo IClientRepo, approvalRepo IApprovalRepo) *AuthService {
	return &AuthService{
		ClientRepo:   clientRepo,
		ApprovalRepo: approvalRepo,
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
	for _, scope := range r.Scopes {
		if !slices.Contains(AllScopes, scope) {
			return ErrInvalidScope
		}
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

var (
	ErrInvalidRequestType = errors.New("invalid request type")
	ErrInvalidClientID    = errors.New("invalid client id. client not found")
	ErrInvalidRedirectURI = errors.New("invalid redirect uri")
	ErrInvalidScope       = errors.New("invalid scope")
)
