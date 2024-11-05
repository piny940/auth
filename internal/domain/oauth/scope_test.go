package oauth

import "testing"

func TestValidScopes(t *testing.T) {
	suites := []struct {
		name     string
		scopes   []TypeScope
		expected error
	}{
		{"valid", []TypeScope{ScopeOpenID}, nil},
		{"invalid", []TypeScope{"invalid"}, ErrInvalidScope},
		{"partially invalid", []TypeScope{ScopeOpenID, "invalid"}, ErrInvalidScope},
		{"empty", []TypeScope{}, nil},
	}
	for _, s := range suites {
		t.Run(s.name, func(t *testing.T) {
			err := ValidScopes(s.scopes)
			if err != s.expected {
				t.Errorf("expected: %v, got: %v", s.expected, err)
			}
		})
	}
}
