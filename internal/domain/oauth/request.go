package oauth

import (
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
	ClientID     string
	RedirectURI  string
	Scopes       []TypeScope
	State        string

	ClientRepo ClientRepo
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
