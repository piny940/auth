package e2e

import (
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure/gateway"
	"auth/internal/infrastructure/model"
	"fmt"
	"io"
	"net/url"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestAuthorizeCodeNotAuthenticated(t *testing.T) {
	const userID = 43234
	const username = "user1"
	const password = "password"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	const clientOwnerID = 32478
	const client1ID = "client1"
	const client2ID = "client2"
	const client3ID = "client3"
	initialUsers := []*model.User{
		{ID: userID, Name: username, EncryptedPassword: string(hashed)},
		{ID: clientOwnerID, Name: "client owner", EncryptedPassword: string(hashed)},
	}
	initialClients := []*model.Client{
		{ID: client1ID, Name: "approved", EncryptedSecret: "secret", UserID: clientOwnerID},
		{ID: client2ID, Name: "approved", EncryptedSecret: "secret", UserID: clientOwnerID},
		{ID: client3ID, Name: "approved", EncryptedSecret: "secret", UserID: clientOwnerID},
	}
	const validRedirectURI = "https://example.com/callback"
	initialRedirectURIs := []*model.RedirectURI{
		{ClientID: client1ID, URI: validRedirectURI},
		{ClientID: client2ID, URI: validRedirectURI},
		{ClientID: client3ID, URI: validRedirectURI},
	}
	validScope := url.QueryEscape(fmt.Sprintf("%s %s", oauth.ScopeOpenID, oauth.ScopeEmail))

	const approval1ID = 32284
	const approval2ID = 32285
	initialApprovals := []*model.Approval{
		{ID: approval1ID, UserID: userID, ClientID: client1ID},
		{ID: approval2ID, UserID: userID, ClientID: client2ID},
	}
	initialApprovalScopes := []*model.ApprovalScope{
		{ApprovalID: approval1ID, ScopeID: gateway.ScopeMapReverse[oauth.ScopeOpenID]},
		{ApprovalID: approval1ID, ScopeID: gateway.ScopeMapReverse[oauth.ScopeEmail]},
		{ApprovalID: approval2ID, ScopeID: gateway.ScopeMapReverse[oauth.ScopeOpenID]},
	}
	const state = "rfejafewiofjwefiojwoefwjofprwjfrawo"

	// serverUrl := os.Getenv("SERVER_URL")
	apiLoginUrl := os.Getenv("API_LOGIN_URL")
	apiApproveUrl := os.Getenv("API_APPROVE_URL")

	suites := []struct {
		name           string
		authenticated  bool
		clientID       string
		RedirectURI    string
		Scope          string
		State          string
		ExpectedStatus int
		Location       *string
	}{
		{"not authenticated", false, client1ID, validRedirectURI, validScope, state, 302, ptr(apiLoginUrl)},
		{"client not found", true, "invalid", validRedirectURI, validScope, state, 400, nil},
		{"invalid redirect uri", true, client1ID, "https://example.com/invalid", validScope, state, 400, nil},
		{"invalid scope", true, client1ID, validRedirectURI, "invalid", state, 400, nil},
		{"partially invalid scope", true, client1ID, validRedirectURI, "openid invalid", state, 400, nil},
		{"empty state", true, client1ID, validRedirectURI, validScope, "", 200, nil},
		{"client not approved", true, client3ID, validRedirectURI, validScope, state, 302, ptr(apiApproveUrl)},
		{"a scope not approved", true, client2ID, validRedirectURI, validScope, state, 302, ptr(apiApproveUrl)},
		{"client approved", true, client1ID, validRedirectURI, validScope, state, 200, nil},
		{"more approved", true, client1ID, validRedirectURI, string(oauth.ScopeOpenID), state, 200, nil},
	}

	for _, suit := range suites {
		t.Run(suit.name, func(t *testing.T) {
			s := newServer(t)
			defer s.Close()
			seed(t, initialUsers, initialClients, initialRedirectURIs, initialApprovals, initialApprovalScopes)

			var cookie *string
			if suit.authenticated {
				_, c := login(t, s, username, password)
				cookie = &c
			}
			query := map[string]string{
				"response_type": "code",
				"client_id":     suit.clientID,
				"redirect_uri":  suit.RedirectURI,
				"scope":         suit.Scope,
			}
			if suit.State != "" {
				query["state"] = suit.State
			}
			res := authedGet(t, s.URL+"/oauth/authorize?"+mapToQuery(t, query), cookie)
			defer res.Body.Close()

			if res.StatusCode != suit.ExpectedStatus {
				body, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				t.Fatalf("expected status code: %d, but got %d. response body: %v", suit.ExpectedStatus, res.StatusCode, string(body))
			}
		})
	}

}
