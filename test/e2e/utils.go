package e2e

import (
	"auth/internal/api"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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
func login(t *testing.T, s *httptest.Server) *api.User {
	t.Helper()

	name := randomString(t, 10)
	password := randomString(t, 16)

	input := api.SessionLoginReq{
		Name:     name,
		Password: password,
	}
	body := toJSON(t, input)
	resp, err := http.Post(s.URL+"/session", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to post: %v", err)
	}
	defer resp.Body.Close()

	resp, err = http.Get(s.URL + "/session")
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}
	defer resp.Body.Close()

	user := &api.User{}
	fromJSONBody(t, resp.Body, user)
	return user
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
