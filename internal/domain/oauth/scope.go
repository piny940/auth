package oauth

import (
	"errors"
	"slices"
)

type TypeScope string

const (
	ScopeOpenID TypeScope = "openid"
	ScopeEmail  TypeScope = "email"
)

var AllScopes = []TypeScope{
	ScopeOpenID,
	ScopeEmail,
}

func ValidScopes(scopes []TypeScope) error {
	for _, s := range scopes {
		if !slices.Contains(AllScopes, s) {
			return ErrInvalidScope
		}
	}
	return nil
}

var ErrInvalidScope = errors.New("invalid scope")
