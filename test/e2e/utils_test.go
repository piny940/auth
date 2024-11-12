package e2e

import (
	"auth/internal/api"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func toJSON(t *testing.T, v interface{}) []byte {
	t.Helper()

	result, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	return result
}
func fromJSONBody(t *testing.T, body io.Reader, v interface{}) {
	t.Helper()

	decoder := json.NewDecoder(body)
	if err := decoder.Decode(v); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
}
func login(t *testing.T, s *httptest.Server, name, password string) (*api.User, string) {
	t.Helper()

	input := &api.SessionLoginReq{
		Name:     name,
		Password: password,
	}
	body := toJSON(t, input)
	resp, err := http.Post(s.URL+"/session", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to post: %v", err)
	}
	defer resp.Body.Close()
	cookie := resp.Header.Get("set-cookie")
	if cookie == "" {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("failed to set cookie: %v", string(body))
	}

	req, err := http.NewRequest(http.MethodGet, s.URL+"/session", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Cookie", cookie)
	resp, err = (&http.Client{}).Do(req)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to get: %v", resp.StatusCode)
	}
	resBody := &struct{ User *api.User }{}
	fromJSONBody(t, resp.Body, resBody)
	return resBody.User, cookie
}

func randomString(t *testing.T, l int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		t.Fatalf("failed to read random: %v", err)
	}
	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}
	return result
}

func seed(
	t *testing.T,
	users []*model.User,
	clients []*model.Client,
	redirectURIs []*model.RedirectURI,
	approvals []*model.Approval,
	approvalScopes []*model.ApprovalScope,
	authCodes []*model.AuthCode,
	authCodeScopes []*model.AuthCodeScope,
) {
	t.Helper()

	db := infrastructure.GetDB()
	query := query.Use(db.Client)

	if users != nil {
		if err := query.User.CreateInBatches(users, len(users)); err != nil {
			t.Fatalf("failed to create users: %v", err)
		}
	}
	if clients != nil {
		if err := query.Client.CreateInBatches(clients, len(clients)); err != nil {
			t.Fatalf("failed to create clients: %v", err)
		}
	}
	if redirectURIs != nil {
		if err := query.RedirectURI.CreateInBatches(redirectURIs, len(redirectURIs)); err != nil {
			t.Fatalf("failed to create redirectURIs: %v", err)
		}
	}
	if approvals != nil {
		if err := query.Approval.CreateInBatches(approvals, len(approvals)); err != nil {
			t.Fatalf("failed to create approvals: %v", err)
		}
	}
	if approvalScopes != nil {
		if err := query.ApprovalScope.CreateInBatches(approvalScopes, len(approvalScopes)); err != nil {
			t.Fatalf("failed to create approvalScopes: %v", err)
		}
	}
	if authCodes != nil {
		if err := query.AuthCode.CreateInBatches(authCodes, len(authCodes)); err != nil {
			t.Fatalf("failed to create authCodes: %v", err)
		}
	}
	if authCodeScopes != nil {
		if err := query.AuthCodeScope.CreateInBatches(authCodeScopes, len(authCodeScopes)); err != nil {
			t.Fatalf("failed to create authCodeScopes: %v", err)
		}
	}
}

func authedGet(t *testing.T, url string, cookie *string) *http.Response {
	t.Helper()

	fmt.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	if cookie != nil {
		req.Header.Set("Cookie", *cookie)
	}
	resp, err := (&http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println("check redirect: ", req.URL, "via: ", via[0].URL)
			return http.ErrUseLastResponse
		},
	}).Do(req)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}
	return resp
}
func authedPost(t *testing.T, url string, cookie *string, body interface{}) *http.Response {
	t.Helper()

	b := toJSON(t, body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	if cookie != nil {
		req.Header.Set("Cookie", *cookie)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}).Do(req)
	if err != nil {
		t.Fatalf("failed to post: %v", err)
	}
	return resp
}

func mapToQuery(t *testing.T, m map[string]string) string {
	t.Helper()

	kvs := make([]string, 0, len(m))
	for k, v := range m {
		kvs = append(kvs, k+"="+url.QueryEscape(v))
	}
	return strings.Join(kvs, "&")
}

func ptr[T any](v T) *T {
	return &v
}
