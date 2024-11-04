package oauth

import (
	"auth/internal/domain"
	"errors"
	"testing"
)

func TestApproved(t *testing.T) {
	approvalRepo := &approvalRepo{
		Approvals: []*Approval{
			{
				ID:       1,
				ClientID: "client1",
				UserID:   1,
				Scopes:   []TypeScope{ScopeOpenID},
			},
		},
	}
	suites := []struct {
		name     string
		clientID ClientID
		userID   domain.UserID
		scopes   []TypeScope
		expected bool
	}{
		{"approved", "client1", 1, []TypeScope{ScopeOpenID}, true},
		{"scope small", "client1", 1, []TypeScope{}, true},
		{"client not match", "client2", 1, []TypeScope{ScopeOpenID}, false},
		{"user not match", "client1", 2, []TypeScope{ScopeOpenID}, false},
		{"scope not match", "client1", 1, []TypeScope{ScopeEmail}, false},
		{"scope too large", "client1", 1, []TypeScope{ScopeOpenID, ScopeEmail}, false},
	}
	for _, s := range suites {
		t.Run(s.name, func(t *testing.T) {
			service := NewApprovalService(approvalRepo)
			result, err := service.Approved(s.clientID, s.userID, s.scopes)
			if err != nil {
				t.Fatal(err)
			}
			if result != s.expected {
				t.Errorf("expected: %t, got: %t", s.expected, result)
			}
		})
	}
}

func TestApprove(t *testing.T) {
	approvalRepo := &approvalRepo{Approvals: []*Approval{}}
	suites := []struct {
		name     string
		clientID ClientID
		userID   domain.UserID
		scopes   []TypeScope
		expected error
	}{
		{"approve", "client1", 1, []TypeScope{ScopeOpenID}, nil},
	}
	for _, suit := range suites {
		t.Run(suit.name, func(t *testing.T) {
			svc := NewApprovalService(approvalRepo)
			err := svc.Approve(suit.clientID, suit.userID, suit.scopes)
			if !errors.Is(err, suit.expected) {
				t.Errorf("expected: %v, got: %v", suit.expected, err)
			}
		})
	}
}