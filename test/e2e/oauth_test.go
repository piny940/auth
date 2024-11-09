package e2e

import (
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure/gateway"
	"auth/internal/infrastructure/model"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/lestrrat-go/jwx/jwk"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthorizeCode(t *testing.T) {
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
		{ID: client2ID, Name: "partially approved", EncryptedSecret: "secret", UserID: clientOwnerID},
		{ID: client3ID, Name: "not approved", EncryptedSecret: "secret", UserID: clientOwnerID},
	}
	const validRedirectURI = "https://example.com/callback"
	initialRedirectURIs := []*model.RedirectURI{
		{ClientID: client1ID, URI: validRedirectURI},
		{ClientID: client2ID, URI: validRedirectURI},
		{ClientID: client3ID, URI: validRedirectURI},
	}
	validScope := "openid email"

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

	serverUrl := os.Getenv("SERVER_URL")
	apiLoginUrl := os.Getenv("API_LOGIN_URL")
	apiApproveUrl := os.Getenv("API_APPROVE_URL")

	suites := []struct {
		name           string
		authenticated  bool
		clientID       string
		redirectURI    string
		scope          string
		state          string
		expectedStatus int
		authCodeIssued bool
		redirectedTo   *string
	}{
		{"not authenticated", false, client1ID, validRedirectURI, validScope, state, 302, false, ptr(apiLoginUrl)},
		{"client not found", true, "invalid", validRedirectURI, validScope, state, 400, false, nil},
		{"invalid redirect uri", true, client1ID, "https://example.com/invalid", validScope, state, 400, false, nil},
		{"invalid scope", true, client1ID, validRedirectURI, "invalid", state, 400, false, nil},
		{"partially invalid scope", true, client1ID, validRedirectURI, "openid invalid", state, 400, false, nil},
		{"empty state", true, client1ID, validRedirectURI, validScope, "", 302, true, ptr(validRedirectURI)},
		{"client not approved", true, client3ID, validRedirectURI, validScope, state, 302, false, ptr(apiApproveUrl)},
		{"a scope not approved", true, client2ID, validRedirectURI, validScope, state, 302, false, ptr(apiApproveUrl)},
		{"client approved", true, client1ID, validRedirectURI, validScope, state, 302, true, ptr(validRedirectURI)},
		{"more approved", true, client1ID, validRedirectURI, string(oauth.ScopeOpenID), state, 302, true, ptr(validRedirectURI)},
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
				"redirect_uri":  suit.redirectURI,
				"scope":         suit.scope,
			}
			if suit.state != "" {
				query["state"] = suit.state
			}
			res := authedGet(t, s.URL+"/oauth/authorize?"+mapToQuery(t, query), cookie)
			defer res.Body.Close()

			if res.StatusCode != suit.expectedStatus {
				body, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				t.Fatalf("expected status code: %d, but got %d. response body: %v", suit.expectedStatus, res.StatusCode, string(body))
			}
			if suit.expectedStatus == 302 {
				actual, err := url.Parse(res.Header.Get("Location"))
				if err != nil {
					t.Fatalf("failed to parse url: %v", err)
				}
				{ // test the redirect target is correct
					actualUrl := fmt.Sprintf("%s://%s%s", actual.Scheme, actual.Host, actual.Path)
					if actualUrl != *suit.redirectedTo {
						t.Errorf("unexpected url: %v, expect %v", actualUrl, *suit.redirectedTo)
					}
				}
				actualQuery := actual.Query()
				if suit.state != "" {
					if actualQuery.Get("state") != suit.state {
						t.Errorf("unexpected state: %v", actualQuery.Get("state"))
					}
				}
				if suit.authCodeIssued {
					if actualQuery.Get("code") == "" {
						t.Errorf("code is not issued")
					}
				} else {
					next, err := url.Parse(actualQuery.Get("next"))
					if err != nil {
						t.Fatalf("failed to parse next url: %v", err)
					}
					actualUrl := fmt.Sprintf("%s://%s%s", next.Scheme, next.Host, next.Path)
					expectedUrl := fmt.Sprintf("%s%s", serverUrl, "/oauth/authorize")
					if actualUrl != expectedUrl {
						t.Errorf("unexpected url: %v, expect %v", actualUrl, expectedUrl)
					}
					nextQuery := next.Query()
					if nextQuery.Get("response_type") != "code" {
						t.Errorf("unexpected response_type: %v, expect code", nextQuery.Get("response_type"))
					}
					if nextQuery.Get("client_id") != suit.clientID {
						t.Errorf("unexpected client_id: %v, expect %v", nextQuery.Get("client_id"), suit.clientID)
					}
					if nextQuery.Get("redirect_uri") != suit.redirectURI {
						t.Errorf("unexpected redirect_uri: %v, expect %v", nextQuery.Get("redirect_uri"), suit.redirectURI)
					}
					if nextQuery.Get("scope") != suit.scope {
						t.Errorf("unexpected scope: %v, expect %v", nextQuery.Get("scope"), suit.scope)
					}
					if suit.state != "" {
						if nextQuery.Get("state") != suit.state {
							t.Errorf("unexpected state: %v, expect %v", nextQuery.Get("state"), suit.state)
						}
					}
				}
			}
		})
	}
}

func TestJwks(t *testing.T) {
	s := newServer(t)
	defer s.Close()

	res, err := http.Get(s.URL + "/oauth/jwks")
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("expected status code: 200, but got %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	_, err = jwk.Parse(body)
	if err != nil {
		t.Fatalf("failed to parse jwks: %v", err)
	}
}
