package e2e

import (
	"auth/internal/api"
	"bytes"
	"net/http"
	"testing"
)

func TestSignupLogin(t *testing.T) {
	s := newServer(t)
	defer s.Close()

	name := randomString(t, 10)
	password := randomString(t, 16)

	{ // signup
		signupInput := &api.UsersReqSignup{
			Name:                 name,
			Password:             password,
			PasswordConfirmation: password,
		}
		body := toJSON(t, signupInput)
		resp, err := http.Post(s.URL+"/users/signup", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("failed to post: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("failed to create user: %v", resp.StatusCode)
		}
	}

	user, cookie := login(t, s, name, password)
	if user.Name != name {
		t.Fatalf("unexpected name: %v", user.Name)
	}

	// test cookie is valid

	req, err := http.NewRequest(http.MethodGet, s.URL+"/session", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Cookie", cookie)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to get: %v", resp.StatusCode)
	}
	resBody := &struct{ User *api.User }{}
	fromJSONBody(t, resp.Body, resBody)
	if resBody.User.Name != name {
		t.Fatalf("unexpected name: %v", resBody.User.Name)
	}
}
