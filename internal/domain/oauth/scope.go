package oauth

import "slices"

type TypeScope string

const (
	ScopeOpenID TypeScope = "openid"
)

var AllScopes = []TypeScope{
	ScopeOpenID,
}

func ValidScopes(scopes []TypeScope) error {
	for _, s := range scopes {
		if !slices.Contains(AllScopes, s) {
			return ErrInvalidScope
		}
	}
	return nil
}
