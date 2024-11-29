package e2e

import (
	"auth/internal/api"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure/gateway"
	"auth/internal/infrastructure/model"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
	"golang.org/x/crypto/bcrypt"
)

func TestOAuthAuthorize(t *testing.T) {
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
		{ID: userID, Email: "test1@example.com", Name: username, EncryptedPassword: string(hashed)},
		{ID: clientOwnerID, Email: "test2@example.com", Name: "client owner", EncryptedPassword: string(hashed)},
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
			seed(t, initialUsers, initialClients, initialRedirectURIs, initialApprovals, initialApprovalScopes, nil, nil)

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

func TestOAuthToken(t *testing.T) {
	const userID = 43234
	const username = "user1"
	const clientOwnerID = 32478
	const password = "password"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	initialUsers := []*model.User{
		{ID: userID, Name: username, Email: "test1@example.com", EncryptedPassword: string(hashed)},
		{ID: clientOwnerID, Name: "client owner", Email: "test2@example.com", EncryptedPassword: "password"},
	}
	const client1ID = "client1"
	const client2ID = "client2"
	const clientSecret = "secret"
	hashed, err = bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	initialClients := []*model.Client{
		{ID: client1ID, Name: "approved", EncryptedSecret: string(hashed), UserID: clientOwnerID},
		{ID: client2ID, Name: "auth code not issued", EncryptedSecret: string(hashed), UserID: clientOwnerID},
	}
	authCodeIds := []int64{3423, 32424, 32423, 2313, 3245}
	const validCode = "validCode"
	const withoutOpenid = "withoutOpenid"
	const emptyScope = "emptyScope"
	const usedCode = "usedCode"
	const expiredCode = "expiredCode"
	const validRedirectURI = "https://example.com/callback"
	expiresAt := time.Now().Add(time.Hour)
	initialAuthCodes := []*model.AuthCode{
		{ID: authCodeIds[0], Value: validCode, RedirectURI: validRedirectURI, Used: false, ExpiresAt: expiresAt, ClientID: client1ID, UserID: userID},
		{ID: authCodeIds[1], Value: usedCode, RedirectURI: validRedirectURI, Used: true, ExpiresAt: expiresAt, ClientID: client1ID, UserID: userID},
		{ID: authCodeIds[2], Value: expiredCode, RedirectURI: validRedirectURI, Used: false, ExpiresAt: time.Now().Add(-time.Hour), ClientID: client1ID, UserID: userID},
		{ID: authCodeIds[3], Value: withoutOpenid, RedirectURI: validRedirectURI, Used: false, ExpiresAt: expiresAt, ClientID: client2ID, UserID: userID},
		{ID: authCodeIds[4], Value: emptyScope, RedirectURI: validRedirectURI, Used: false, ExpiresAt: expiresAt, ClientID: client1ID, UserID: userID},
	}
	initialAuthCodeScopes := []*model.AuthCodeScope{
		{AuthCodeID: authCodeIds[0], ScopeID: gateway.ScopeMapReverse[oauth.ScopeOpenID]},
		{AuthCodeID: authCodeIds[0], ScopeID: gateway.ScopeMapReverse[oauth.ScopeEmail]},
		{AuthCodeID: authCodeIds[3], ScopeID: gateway.ScopeMapReverse[oauth.ScopeEmail]},
	}

	suites := []struct {
		name              string
		grantType         api.OAuthTokenGrantType
		code              string
		redirectURI       string
		clientId          string
		clientSecret      string
		expected          int
		idTokenIssued     bool
		accessTokenIssued bool
	}{
		{"invalid grant type", "invalid", validCode, validRedirectURI, client1ID, clientSecret, 400, false, false},
		{"code not found", api.AuthorizationCode, "invalid", validRedirectURI, client1ID, clientSecret, 400, false, false},
		{"code used", api.AuthorizationCode, usedCode, validRedirectURI, client1ID, clientSecret, 400, false, false},
		{"code expired", api.AuthorizationCode, expiredCode, validRedirectURI, client1ID, clientSecret, 400, false, false},
		{"redirect uri not match", api.AuthorizationCode, validCode, "https://example.com/invalid", client1ID, clientSecret, 400, false, false},
		{"client not found", api.AuthorizationCode, validCode, validRedirectURI, "invalid", clientSecret, 403, false, false},
		{"client secret not match", api.AuthorizationCode, validCode, validRedirectURI, client1ID, "invalid", 403, false, false},
		{"success", api.AuthorizationCode, validCode, validRedirectURI, client1ID, clientSecret, 200, true, true},
		{"without openid", api.AuthorizationCode, withoutOpenid, validRedirectURI, client2ID, clientSecret, 200, false, true},
		{"empty scope", api.AuthorizationCode, emptyScope, validRedirectURI, client1ID, clientSecret, 200, false, true},
	}

	for _, suit := range suites {
		t.Run(suit.name, func(t *testing.T) {
			s := newServer(t)
			defer s.Close()

			seed(t, initialUsers, initialClients, nil, nil, nil, initialAuthCodes, initialAuthCodeScopes)

			form := url.Values{}
			form.Add("grant_type", string(suit.grantType))
			form.Add("code", suit.code)
			form.Add("redirect_uri", suit.redirectURI)
			form.Add("client_id", suit.clientId)
			reqBody := strings.NewReader(form.Encode())
			req, err := http.NewRequest(http.MethodPost, s.URL+"/oauth/token", reqBody)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.SetBasicAuth(suit.clientId, suit.clientSecret)

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != suit.expected {
				body, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				t.Fatalf("expected status code: %d, but got %d. response body: %v", suit.expected, res.StatusCode, string(body))
			}
			if suit.expected != 200 {
				return
			}
			data := &api.OAuthTokenRes{}
			fromJSONBody(t, res.Body, data)
			if suit.idTokenIssued {
				if data.IdToken == nil {
					t.Errorf("id token is not issued")
				}
				{ // validate id token
					claims := parseToken(t, s, *data.IdToken).(jwt.MapClaims)
					if claims["iss"] != os.Getenv("OAUTH_ISSUER") {
						t.Errorf("unexpected aud: %v", claims["aud"])
					}
					const gap = 5 * time.Second
					expiresAt := time.Unix(int64(claims["exp"].(float64)), 0)
					expectedExpiresAt := time.Now().Add(oauth.ID_TOKEN_TTL)
					if expiresAt.Before(expectedExpiresAt.Add(-gap)) || expiresAt.After(expectedExpiresAt.Add(gap)) {
						t.Errorf("unexpected exp: %v", expiresAt)
					}
					if claims["sub"] != fmt.Sprintf("id:%d;name:%s", userID, username) {
						t.Errorf("unexpected sub: %v", claims["sub"])
					}
					issuedAt := time.Unix(int64(claims["iat"].(float64)), 0)
					if issuedAt.Before(time.Now().Add(-gap)) || issuedAt.After(time.Now().Add(gap)) {
						t.Errorf("unexpected iat: %v", issuedAt)
					}
					if claims["jti"] == "" {
						t.Errorf("jti is empty")
					}
				}
			}
			if suit.accessTokenIssued {
				if data.AccessToken == "" {
					t.Errorf("access token is not issued")
				}
			}
			{ // validate access token
				if data.AccessToken == "" {
					t.Fatalf("access token is empty")
				}
				claims := parseToken(t, s, data.AccessToken).(jwt.MapClaims)
				if claims["iss"] != os.Getenv("OAUTH_ISSUER") {
					t.Errorf("unexpected aud: %v", claims["aud"])
				}
				const gap = 5 * time.Second
				expiresAt := time.Unix(int64(claims["exp"].(float64)), 0)
				expectedExpiresAt := time.Now().Add(oauth.ACCESS_TOKEN_TTL)
				if expiresAt.Before(expectedExpiresAt.Add(-gap)) || expiresAt.After(expectedExpiresAt.Add(gap)) {
					t.Errorf("unexpected exp: %v", expiresAt)
				}
				if claims["sub"] != fmt.Sprintf("id:%d;name:%s", userID, username) {
					t.Errorf("unexpected sub: %v", claims["sub"])
				}
				issuedAt := time.Unix(int64(claims["iat"].(float64)), 0)
				if issuedAt.Before(time.Now().Add(-gap)) || issuedAt.After(time.Now().Add(gap)) {
					t.Errorf("unexpected iat: %v", issuedAt)
				}
				if claims["jti"] == "" {
					t.Errorf("jti is empty")
				}
			}

			// test auth code can be used only once
			res, err = http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			if res.StatusCode != 400 {
				t.Fatalf("expected status code: 400, but got %d", res.StatusCode)
			}
		})
	}
}

func TestAuthTime(t *testing.T) {
	const userID = 43234
	const username = "user1"
	const password = "password"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	const clientOwnerID = 32478
	initialUsers := []*model.User{
		{ID: clientOwnerID, Name: "client owner", Email: "test1@example.com", EncryptedPassword: "password"},
		{ID: userID, Name: username, Email: "test2@example.com", EncryptedPassword: string(hashed)},
	}
	const clientID = "client1"
	const clientSecret = "secret"
	hashed, err = bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	initialClients := []*model.Client{
		{ID: clientID, Name: "client1", EncryptedSecret: string(hashed), UserID: clientOwnerID},
	}
	const redirectURI = "https://example.com/callback"
	initialRedirectURIs := []*model.RedirectURI{
		{ClientID: "client1", URI: redirectURI},
	}

	suites := []struct {
		name         string
		maxAge       int
		expectAuthOk bool
	}{
		{"auth too old", 1, false},
		{"auth ok", 2, true},
	}
	for _, suit := range suites {
		t.Run(suit.name, func(t *testing.T) {
			s := newServer(t)
			defer s.Close()
			seed(t, initialUsers, initialClients, initialRedirectURIs, nil, nil, nil, nil)

			authTime := time.Now()
			_, c := login(t, s, username, password)

			time.Sleep(1 * time.Second)

			approveReq := &api.ApprovalsApproveReq{
				ClientId: clientID,
				Scope:    string(oauth.ScopeOpenID),
			}
			approveRes := authedPost(t, s.URL+"/account/approvals", &c, approveReq)
			defer approveRes.Body.Close()
			if approveRes.StatusCode != 204 {
				t.Fatalf("expected status code: 204, but got %d", approveRes.StatusCode)
			}

			query := map[string]string{
				"response_type": "code",
				"client_id":     clientID,
				"redirect_uri":  redirectURI,
				"scope":         string(oauth.ScopeOpenID),
				"max_age":       strconv.Itoa(suit.maxAge),
			}
			authRes := authedGet(t, s.URL+"/oauth/authorize?"+mapToQuery(t, query), &c)
			defer authRes.Body.Close()
			if authRes.StatusCode != 302 {
				t.Fatalf("expected status code: 302, but got %d", authRes.StatusCode)
			}
			redirected, err := url.Parse(authRes.Header.Get("Location"))
			if err != nil {
				t.Fatalf("failed to parse url: %v", err)
			}
			if !suit.expectAuthOk {
				if redirected.Query().Get("code") != "" {
					t.Errorf("code is not expected to be issued")
				}
				if redirected.Query().Get("error") != string(api.OAuthAuthorizeErrUnauthorizedClient) {
					t.Errorf("unexpected error: %v", redirected.Query().Get("error"))
				}
				return
			}
			code := redirected.Query().Get("code")
			if code == "" {
				t.Fatalf("code is empty")
			}

			form := url.Values{}
			form.Add("grant_type", string(api.AuthorizationCode))
			form.Add("code", code)
			form.Add("redirect_uri", redirectURI)
			form.Add("client_id", clientID)
			reqBody := strings.NewReader(form.Encode())
			req, err := http.NewRequest(http.MethodPost, s.URL+"/oauth/token", reqBody)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.SetBasicAuth(clientID, clientSecret)

			tokenRes, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer tokenRes.Body.Close()
			if tokenRes.StatusCode != 200 {
				t.Fatalf("expected status code: 200, but got %d", tokenRes.StatusCode)
			}
			data := &api.OAuthTokenRes{}
			fromJSONBody(t, tokenRes.Body, data)
			{ // test id token's auth_time
				if data.IdToken == nil {
					t.Fatalf("id token is empty")
				}
				claims := parseToken(t, s, *data.IdToken).(jwt.MapClaims)
				if claims["auth_time"] == nil {
					t.Fatalf("auth_time is empty")
				}
				actual := time.Unix(int64(claims["auth_time"].(float64)), 0)
				gap := 1 * time.Second
				if actual.Before(authTime.Add(-gap)) || actual.After(authTime.Add(gap)) {
					t.Errorf("unexpected auth_time: %v, expect %v", actual, authTime)
				}
			}
			{ // test access token's auth_time
				if data.AccessToken == "" {
					t.Fatalf("access token is empty")
				}
				claims := parseToken(t, s, data.AccessToken).(jwt.MapClaims)
				if claims["auth_time"] == nil {
					t.Fatalf("auth_time is empty")
				}
				actual := time.Unix(int64(claims["auth_time"].(float64)), 0)
				gap := 1 * time.Second
				if actual.Before(authTime.Add(-gap)) || actual.After(authTime.Add(gap)) {
					t.Errorf("unexpected auth_time: %v, expect %v", actual, authTime)
				}
			}
		})
	}
}

func TestOAuthJwks(t *testing.T) {
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

func parseToken(t *testing.T, s *httptest.Server, token string) jwt.Claims {
	t.Helper()

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
	set, err := jwk.Parse(body)
	if err != nil {
		t.Fatalf("failed to parse jwks: %v", err)
	}
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid")
		}
		key, ok := set.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key not found")
		}
		var pubKey interface{}
		if err := key.Raw(&pubKey); err != nil {
			return nil, fmt.Errorf("failed to get public key: %w", err)
		}
		return pubKey, nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	return tokenObj.Claims
}
