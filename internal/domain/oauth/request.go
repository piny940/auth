package oauth

import (
	"auth/internal/domain"
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

	ClientRepo   ClientRepo
	ApprovalRepo IApprovalRepo
}

type ApprovalID int64
type Approval struct {
	ID        ApprovalID
	ClientID  ClientID
	UserID    int64
	Scopes    []TypeScope
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IApprovalRepo interface {
	Find(clientID ClientID, userID domain.UserID) (*Approval, error)
}

func (r *AuthRequest) Validate() error {
	if !slices.Contains(AllResponseTypes, r.ResponseType) {
		return ErrInvalidRequestType{}
	}
	client, err := r.ClientRepo.FindByID(r.ClientID)
	if err != nil {
		return err
	}
	if !slices.Contains(client.RedirectURIs, r.RedirectURI) {
		return ErrInvalidRedirectURI{}
	}
	for _, scope := range r.Scopes {
		if !slices.Contains(AllScopes, scope) {
			return ErrInvalidScope{}
		}
	}

	return nil
}

func (r *AuthRequest) ApprovedBy(user *domain.User) (bool, error) {
	approval, err := r.ApprovalRepo.Find(r.ClientID, user.ID)
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

type ErrInvalidRequestType struct{}

func (e ErrInvalidRequestType) Error() string {
	return "invalid request type"
}

type ErrInvalidRedirectURI struct{}

func (e ErrInvalidRedirectURI) Error() string {
	return "invalid redirect uri"
}

type ErrInvalidScope struct{}

func (e ErrInvalidScope) Error() string {
	return "invalid scope"
}
