package oauth

import (
	"auth/internal/domain"
	"slices"
	"testing"
	"time"
)

type authCodeRepo struct {
	authCodes []*AuthCode
}

var _ IAuthCodeRepo = &authCodeRepo{}

func (a *authCodeRepo) Create(value string, clientID ClientID, userID domain.UserID, scopes []TypeScope, expiresAt time.Time) error {
	a.authCodes = append(a.authCodes, &AuthCode{
		Value:     value,
		ClientID:  clientID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Scopes:    scopes,
	})
	return nil
}

type approvalRepo struct{}

var _ IApprovalRepo = &approvalRepo{}

func (a *approvalRepo) Create(clientID ClientID, userID domain.UserID, scopes []TypeScope) error {
	panic("unimplemented")
}

func (a *approvalRepo) Find(clientID ClientID, userID domain.UserID) (*Approval, error) {
	panic("unimplemented")
}

type clientRepo struct{}

var _ IClientRepo = &clientRepo{}

func (c *clientRepo) FindByID(id ClientID) (*Client, error) {
	panic("unimplemented")
}

func TestAuthCodeServiceIssueAuthCode(t *testing.T) {
	s := &AuthCodeService{
		AuthCodeRepo: &authCodeRepo{authCodes: []*AuthCode{}},
	}
	client := &Client{ID: "client_id"}
	user := &domain.User{ID: 1}
	scopes := []TypeScope{"scope1", "scope2"}
	code1, err := s.IssueAuthCode(client.ID, user.ID, scopes)
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
	code2, err := s.IssueAuthCode(client.ID, user.ID, scopes)
	if err != nil {
		t.Fatal(err)
	}
	if code1.Value == code2.Value {
		t.Errorf("want different value, got same value")
	}
}
