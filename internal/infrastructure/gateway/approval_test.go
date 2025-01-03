package gateway

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"errors"
	"slices"
	"testing"
)

func TestApprovalApprove(t *testing.T) {
	const userID = 87392
	const approvalID = 43289
	initialUsers := []*model.User{
		{ID: userID, Name: "test", EncryptedPassword: "test"},
	}
	initialClients := []*model.Client{
		{ID: "client1", EncryptedSecret: "", UserID: userID, Name: "Not Approved"},
		{ID: "client2", EncryptedSecret: "", UserID: userID, Name: "Partially Approved"},
	}
	initialApprovals := []*model.Approval{
		{ID: approvalID, ClientID: "client2", UserID: userID},
	}
	initialApprovalScopes := []*model.ApprovalScope{
		{ApprovalID: approvalID, ScopeID: ScopeMapReverse[oauth.ScopeOpenID]},
	}
	suites := []struct {
		name           string
		clientID       oauth.ClientID
		userID         domain.UserID
		scopes         []oauth.TypeScope
		expectedScopes []oauth.TypeScope
	}{
		{"fully new approval", "client1", userID, []oauth.TypeScope{oauth.ScopeOpenID, oauth.ScopeEmail}, []oauth.TypeScope{oauth.ScopeOpenID, oauth.ScopeEmail}},
		{"partially new scopes", "client2", userID, []oauth.TypeScope{oauth.ScopeOpenID, oauth.ScopeEmail}, []oauth.TypeScope{oauth.ScopeOpenID, oauth.ScopeEmail}},
		{"new scope", "client2", userID, []oauth.TypeScope{oauth.ScopeEmail}, []oauth.TypeScope{oauth.ScopeOpenID, oauth.ScopeEmail}},
		{"no new scopes", "client2", userID, []oauth.TypeScope{oauth.ScopeOpenID}, []oauth.TypeScope{oauth.ScopeOpenID}},
		{"empty scopes", "client2", userID, []oauth.TypeScope{}, []oauth.TypeScope{oauth.ScopeOpenID}},
	}
	for _, s := range suites {
		t.Run(s.name, func(t *testing.T) {
			setup(t)
			db := infrastructure.GetDB()
			query := query.Use(db.Client)
			if err := query.User.CreateInBatches(initialUsers, len(initialUsers)); err != nil {
				t.Fatal(err)
			}
			if err := query.Client.CreateInBatches(initialClients, len(initialClients)); err != nil {
				t.Fatal(err)
			}
			if err := query.Approval.CreateInBatches(initialApprovals, len(initialApprovals)); err != nil {
				t.Fatal(err)
			}
			if err := query.ApprovalScope.CreateInBatches(initialApprovalScopes, len(initialApprovalScopes)); err != nil {
				t.Fatal(err)
			}
			approvalRepo := NewApprovalRepo(db)
			if err := approvalRepo.Approve(s.clientID, s.userID, s.scopes); err != nil {
				t.Fatal(err)
			}
			actual, err := query.Approval.Where(
				query.Approval.UserID.Eq(int64(s.userID)),
				query.Approval.ClientID.Eq(string(s.clientID)),
			).First()
			if err != nil {
				t.Fatal(err)
			}
			scopes, err := query.ApprovalScope.Where(query.ApprovalScope.ApprovalID.Eq(actual.ID)).Find()
			if err != nil {
				t.Fatal(err)
			}
			if len(scopes) != len(s.expectedScopes) {
				t.Errorf("expected: %v, got: %v", s.expectedScopes, scopes)
			}
			for _, scope := range scopes {
				if !slices.Contains(s.expectedScopes, ScopeMap[scope.ScopeID]) {
					t.Errorf("expected: %v, got: %v", s.expectedScopes, scopes)
				}
			}
		})
	}
}

func TestApprovalFind(t *testing.T) {
	const userID = 87392
	const clientID = "client1"
	initialUsers := []*model.User{
		{ID: userID, Name: "test", EncryptedPassword: "test"},
	}
	initialClients := []*model.Client{
		{ID: clientID, EncryptedSecret: "", UserID: userID, Name: "Not Approved"},
	}
	suites := []struct {
		name     string
		ClientID oauth.ClientID
		UserID   domain.UserID
		Scopes   []oauth.TypeScope
		found    bool
	}{
		{"with two scopes", clientID, userID, []oauth.TypeScope{oauth.ScopeOpenID, oauth.ScopeEmail}, true},
		{"with one scope", clientID, userID, []oauth.TypeScope{oauth.ScopeOpenID}, true},
		{"empty scope", clientID, userID, []oauth.TypeScope{}, true},
		{"not found", clientID, userID, []oauth.TypeScope{}, false},
	}
	for _, s := range suites {
		t.Run(s.name, func(t *testing.T) {
			setup(t)
			db := infrastructure.GetDB()
			query := query.Use(db.Client)
			if err := query.User.CreateInBatches(initialUsers, len(initialUsers)); err != nil {
				t.Fatal(err)
			}
			if err := query.Client.CreateInBatches(initialClients, len(initialClients)); err != nil {
				t.Fatal(err)
			}
			repo := NewApprovalRepo(db)
			if s.found {
				repo.Approve(s.ClientID, s.UserID, s.Scopes)
			}
			actual, err := repo.Find(s.ClientID, s.UserID)
			if !s.found {
				if !errors.Is(err, domain.ErrRecordNotFound) {
					t.Errorf("expected: %v, got: %v", domain.ErrRecordNotFound, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(actual.Scopes) != len(s.Scopes) {
				t.Errorf("expected: %v, got: %v", s.Scopes, actual.Scopes)
			}
			for _, scope := range actual.Scopes {
				if !slices.Contains(s.Scopes, scope) {
					t.Errorf("expected: %v, got: %v", s.Scopes, actual.Scopes)
				}
			}
		})
	}
}
