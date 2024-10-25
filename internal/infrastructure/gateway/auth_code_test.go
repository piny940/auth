package gateway

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"testing"
	"time"
)

func TestAuthCodeCreate(t *testing.T) {
	setup(t)

	db := infrastructure.GetDB()
	query := query.Use(db.Client)
	authCodeRepo := NewAuthCodeRepo(db)
	scopes := []oauth.TypeScope{oauth.ScopeOpenID, oauth.ScopeOpenID}
	expiresAt := time.Now()
	clientID := "test_client_id"
	userID := int64(1)
	query.User.Create(&model.User{ID: userID, Name: "test", EncryptedPassword: "test"})
	query.User.Create(&model.User{ID: 2, Name: "test2", EncryptedPassword: "test2"})
	query.Client.Create(&model.Client{
		ID:              clientID,
		EncryptedSecret: "",
		UserID:          2,
	})
	err := authCodeRepo.Create(
		"test_value",
		oauth.ClientID(clientID),
		domain.UserID(userID),
		scopes,
		expiresAt,
		"test_redirect_uri",
	)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := query.AuthCode.Where(query.AuthCode.Value.Eq("test_value")).First()
	if err != nil {
		t.Fatal(err)
	}
	if actual.ClientID != "test_client_id" {
		t.Errorf("unexpected actual.ClientID: %s", actual.ClientID)
	}
	if actual.UserID != 1 {
		t.Errorf("unexpected actual.UserID: %d", actual.UserID)
	}
	if actual.ExpiresAt.Unix() != expiresAt.Unix() { // DBはnano秒精度を持たないため秒単位で比較する
		t.Errorf("expected ExpiresAt: %s, but got %s", expiresAt, actual.ExpiresAt)
	}
	if actual.RedirectURI != "test_redirect_uri" {
		t.Errorf("unexpected actual.RedirectURI: %s", actual.RedirectURI)
	}
	if actual.Used {
		t.Errorf("unexpected actual.Used: %t", actual.Used)
	}
}
