package oauth

import (
	"auth/internal/domain"
	"slices"
	"testing"
	"time"
)

func TestAuthCodeServiceIssueAuthCode(t *testing.T) {
	s := NewAuthCodeService(&authCodeRepo{authCodes: []*AuthCode{}})
	client := &Client{ID: "client_id"}
	user := &domain.User{ID: 1}
	scopes := []TypeScope{"scope1", "scope2"}
	authTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	code1, err := s.IssueAuthCode(client.ID, user.ID, authTime, scopes, "redirect_uri")
	if err != nil {
		t.Fatal(err)
	}
	if code1.ClientID != client.ID {
		t.Errorf("want %s, got %s", client.ID, code1.ClientID)
	}
	if code1.UserID != user.ID {
		t.Errorf("want %d, got %d", user.ID, code1.UserID)
	}
	if !slices.Equal(scopes, code1.Scopes) {
		t.Errorf("want %v, got %v", scopes, code1.Scopes)
	}
	if code1.ExpiresAt.IsZero() {
		t.Error("want non-zero, got zero")
	}
	if code1.Used {
		t.Error("want false, got true")
	}
	if code1.AuthTime != authTime {
		t.Errorf("want %v, got %v", authTime, code1.AuthTime)
	}
	code2, err := s.IssueAuthCode(client.ID, user.ID, authTime, scopes, "redirect_uri")
	if err != nil {
		t.Fatal(err)
	}
	if code1.Value == code2.Value {
		t.Errorf("want different value, got same value")
	}
}

func TestAuthCodeVerify(t *testing.T) {
	now := time.Now()
	authTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	oauthCodes := []*AuthCode{
		{"active", "client_id", 1, now.Add(time.Minute), false, authTime, "redirect_uri", []TypeScope{"scope1", "scope2"}},
		{"used", "client_id", 1, now.Add(time.Minute), true, authTime, "redirect_uri", []TypeScope{"scope1", "scope2"}},
		{"expired", "client_id", 1, now.Add(-time.Minute), false, authTime, "redirect_uri", []TypeScope{"scope1", "scope2"}},
	}
	suites := []struct {
		Name        string
		Value       string
		ClientID    ClientID
		RedirectURI string
		WantError   bool
	}{
		{"valid", "active", "client_id", "redirect_uri", false},
		{"used", "used", "client_id", "redirect_uri", true},
		{"expired", "expired", "client_id", "redirect_uri", true},
		{"not found", "not_found", "client_id", "redirect_uri", true},
		{"invalid client id", "active", "invalid_client_id", "redirect_uri", true},
		{"invalid redirect url", "active", "client_id", "invalid_redirect_uri", true},
	}
	for _, suit := range suites {
		t.Run(suit.Name, func(t *testing.T) {
			s := &AuthCodeService{
				AuthCodeRepo: &authCodeRepo{authCodes: oauthCodes},
			}
			_, err := s.Verify(suit.Value, suit.ClientID, suit.RedirectURI)
			if suit.WantError {
				if err == nil {
					t.Error("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
