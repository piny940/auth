package e2e

import (
	"auth/internal/api"
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"fmt"
	"net/http"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestUserInfo(t *testing.T) {
	s := newServer(t)
	defer s.Close()

	db := infrastructure.GetDB()
	query := query.Use(db.Client)
	userID := 43829
	name := randomString(t, 10)
	password := randomString(t, 16)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate hash: %v", err)
	}

	query.User.Create(&model.User{
		ID:                int64(userID),
		Name:              name,
		EncryptedPassword: string(hash),
	})
	token := accessToken(t, domain.UserID(userID), []oauth.TypeScope{oauth.ScopeEmail})
	req, err := http.NewRequest(http.MethodGet, s.URL+"/userinfo", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %v", resp.StatusCode)
	}
	resBody := &api.UserinfoUserinfoRes{}
	fromJSONBody(t, resp.Body, resBody)
	if resBody.Sub != fmt.Sprintf("id:%v;name:%v", userID, name) {
		t.Fatalf("unexpected name: %v", resBody.Sub)
	}
}
