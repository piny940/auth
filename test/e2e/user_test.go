package e2e

import (
	"auth/internal/api"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/query"
	"bytes"
	"net/http"
	"testing"
)

func TestSignupLogin(t *testing.T) {
	s := newServer(t)
	defer s.Close()

	db := infrastructure.GetDB()
	query := query.Use(db.Client)

	before, err := query.User.Find()
	if err != nil {
		t.Fatal(err)
	}

	input := api.UsersReqSignup{
		Name:                 "test",
		Password:             "password",
		PasswordConfirmation: "password",
	}
	body := toJSON(t, input)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	resp, err := http.Post(s.URL+"/users/signup", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("failed to create user: %v", resp.StatusCode)
	}
	after, err := query.User.Find()
	if err != nil {
		t.Fatal(err)
	}
	if len(before)+1 != len(after) {
		t.Errorf("failed to create user: before: %v, after: %v", before, after)
	}
	if resp.Header.Get("Set-Cookie") == "" {
		t.Errorf("failed to set cookie")
	}

	loginInput := api.SessionLoginReq{
		Name:     input.Name,
		Password: input.Password,
	}
	loginBody := toJSON(t, loginInput)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	resp, err = http.Post(s.URL+"/session", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		t.Fatalf("failed to post: %v", err)
	}
	defer resp.Body.Close()
	if resp.Header.Get("Set-Cookie") == "" {
		t.Errorf("failed to set cookie")
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("failed to login: %v", resp.StatusCode)
	}
}

func TestSessionMe(t *testing.T) {
	s := newServer(t)
	defer s.Close()

	login(t, s)
}
