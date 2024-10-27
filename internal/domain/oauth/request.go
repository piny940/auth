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

type AuthRequest struct {
	ResponseType TypeResponseType
	ClientID     ClientID
	RedirectURI  string
	Scopes       []TypeScope
	State        *string
}

type RequestService struct {
	ClientRepo   IClientRepo
	ApprovalRepo IApprovalRepo
	AuthCodeRepo IAuthCodeRepo
}

func NewRequestService(clientRepo IClientRepo, approvalRepo IApprovalRepo, authCodeRepo IAuthCodeRepo) *RequestService {
	return &RequestService{
		ClientRepo:   clientRepo,
		ApprovalRepo: approvalRepo,
		AuthCodeRepo: authCodeRepo,
	}
}

func (s *RequestService) Validate(r *AuthRequest) error {
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

var (
	ErrInvalidRequestType = errors.New("invalid request type")
	ErrInvalidClientID    = errors.New("invalid client id. client not found")
	ErrInvalidRedirectURI = errors.New("invalid redirect uri")
	ErrInvalidScope       = errors.New("invalid scope")
)
