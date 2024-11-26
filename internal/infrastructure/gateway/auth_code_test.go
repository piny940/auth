package gateway

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"errors"
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
	authTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	clientID := "test_client_id"
	userID := int64(1)
	query.User.Create(&model.User{ID: userID, Email: "test1@example.com", Name: "test", EncryptedPassword: "test"})
	query.User.Create(&model.User{ID: 2, Email: "test2@example.com", Name: "test2", EncryptedPassword: "test2"})
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
		authTime,
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
	if actual.AuthTime.Unix() != authTime.Unix() {
		t.Errorf("unexpected actual.AuthTime: %s", actual.AuthTime)
	}
	if actual.RedirectURI != "test_redirect_uri" {
		t.Errorf("unexpected actual.RedirectURI: %s", actual.RedirectURI)
	}
	if actual.Used {
		t.Errorf("unexpected actual.Used: %t", actual.Used)
	}
	actualScopes, err := query.AuthCodeScope.Where(
		query.AuthCodeScope.AuthCodeID.Eq(actual.ID),
	).Find()
	if err != nil {
		t.Fatal(err)
	}
	if len(actualScopes) != 1 { // 重複は排除される
		t.Fatalf("unexpected len(actualScopes): %d", len(actualScopes))
	}
	if actualScopes[0].ScopeID != ScopeMapReverse[oauth.ScopeOpenID] {
		t.Errorf("unexpected actualScopes[0].ScopeID: %d", actualScopes[0].ScopeID)
	}
}

func TestAuthCodeFind(t *testing.T) {
	setup(t)

	db := infrastructure.GetDB()
	query := query.Use(db.Client)
	authCodeRepo := NewAuthCodeRepo(db)
	scopes := []oauth.TypeScope{oauth.ScopeOpenID, oauth.ScopeOpenID}
	expiresAt := time.Now()
	authTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	clientID := "test_client_id"
	userID := int64(1)
	value := "test_value"
	query.User.Create(&model.User{ID: userID, Email: "test1@example.com", Name: "test", EncryptedPassword: "test"})
	query.Client.Create(&model.Client{ID: clientID, EncryptedSecret: "", UserID: 1})
	err := authCodeRepo.Create(value, oauth.ClientID(clientID), domain.UserID(userID), scopes, expiresAt, authTime, "test_redirect_uri")
	if err != nil {
		t.Fatal(err)
	}
	authCode, err := authCodeRepo.Find(value)
	if err != nil {
		t.Fatal(err)
	}
	if authCode.ClientID != oauth.ClientID(clientID) {
		t.Errorf("unexpected authCode.ClientID: %s", authCode.ClientID)
	}
	if authCode.UserID != domain.UserID(userID) {
		t.Errorf("unexpected authCode.UserID: %d", authCode.UserID)
	}
	if authCode.ExpiresAt.Unix() != expiresAt.Unix() { // DBはnano秒精度を持たないため秒単位で比較する
		t.Errorf("expected ExpiresAt: %s, but got %s", expiresAt, authCode.ExpiresAt)
	}
	if authCode.AuthTime.Unix() != authTime.Unix() {
		t.Errorf("unexpected authCode.AuthTime: %s", authCode.AuthTime)
	}

	if authCode.RedirectURI != "test_redirect_uri" {
		t.Errorf("unexpected authCode.RedirectURI: %s", authCode.RedirectURI)
	}
	if authCode.Used {
		t.Errorf("unexpected authCode.Used: %t", authCode.Used)
	}
	if len(authCode.Scopes) != 1 {
		t.Fatalf("unexpected len(authCode.Scopes): %d", len(authCode.Scopes))
	}
	if authCode.Scopes[0] != oauth.ScopeOpenID {
		t.Errorf("unexpected authCode.Scopes[0]: %s", authCode.Scopes[0])
	}
}

func TestAuthCodeRepoFindEmpty(t *testing.T) {
	setup(t)

	db := infrastructure.GetDB()
	authCodeRepo := NewAuthCodeRepo(db)
	_, err := authCodeRepo.Find("test")
	if !errors.Is(err, domain.ErrRecordNotFound) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAuthCodeUse(t *testing.T) {
	setup(t)

	db := infrastructure.GetDB()
	query := query.Use(db.Client)
	authCodeRepo := NewAuthCodeRepo(db)
	scopes := []oauth.TypeScope{oauth.ScopeOpenID, oauth.ScopeOpenID}
	expiresAt := time.Now()
	authTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	clientID := "test_client_id"
	userID := int64(1)
	value := "test_value"
	query.User.Create(&model.User{ID: userID, Name: "test", EncryptedPassword: "test"})
	query.Client.Create(&model.Client{ID: clientID, EncryptedSecret: "", UserID: 1})
	err := authCodeRepo.Create(value, oauth.ClientID(clientID), domain.UserID(userID), scopes, expiresAt, authTime, "test_redirect_uri")
	if err != nil {
		t.Fatal(err)
	}
	before, err := authCodeRepo.Find(value)
	if err != nil {
		t.Fatal(err)
	}
	if before.Used {
		t.Errorf("unexpected before.Used: %t", before.Used)
	}
	err = authCodeRepo.Use(value)
	if err != nil {
		t.Fatal(err)
	}
	authCode, err := authCodeRepo.Find(value)
	if err != nil {
		t.Fatal(err)
	}
	if !authCode.Used {
		t.Errorf("unexpected authCode.Used: %t", authCode.Used)
	}
}
