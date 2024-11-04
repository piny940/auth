package gateway

import (
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"slices"
	"testing"
)

func TestClientCreate(t *testing.T) {
	const userID = 45328
	initialUsers := []*model.User{
		{Name: "test1", ID: userID, EncryptedPassword: "password"},
	}
	initialClients := []*model.Client{
		{ID: "client1", Name: "client1", UserID: userID},
	}
	initialRedirectURIs := []*model.RedirectURI{
		{ClientID: "client1", URI: "http://example.com"},
	}
	suites := []struct {
		name         string
		clientID     oauth.ClientID
		redirectURIs []string
	}{
		{"with a URI", "client2", []string{"http://example.com"}},
		{"with multiple redirect URIs", "client2", []string{"http://example.com", "http://example2.com"}},
		{"without a URI", "client2", []string{}},
	}
	for _, s := range suites {
		t.Run(s.name, func(t *testing.T) {
			setup(t)
			db := infrastructure.GetDB()
			query := query.Use(db.Client)
			if err := query.User.CreateInBatches(initialUsers, len(initialUsers)); err != nil {
				t.Fatal(err)
			}
			if err := query.Client.CreateInBatches(initialClients, len(initialClients)); err != nil {
				t.Fatal(err)
			}
			if err := query.RedirectURI.CreateInBatches(initialRedirectURIs, len(initialRedirectURIs)); err != nil {
				t.Fatal(err)
			}

			clientRepo := NewClientRepo(db)
			err := clientRepo.Create(&oauth.ClientInput{
				ID:              s.clientID,
				Name:            string(s.clientID),
				EncryptedSecret: "secret",
				UserID:          userID,
				RedirectURIs:    s.redirectURIs,
			})
			if err != nil {
				t.Fatal(err)
			}
			client, err := query.Client.Where(query.Client.ID.Eq(string(s.clientID))).First()
			if err != nil {
				t.Fatal(err)
			}
			uris, err := query.RedirectURI.Where(query.RedirectURI.ClientID.Eq(client.ID)).Find()
			if err != nil {
				t.Fatal(err)
			}
			if len(uris) != len(s.redirectURIs) {
				t.Fatalf("expected %d redirect URIs, got %d", len(s.redirectURIs), len(uris))
			}
			for _, uri := range uris {
				if !slices.Contains(s.redirectURIs, uri.URI) {
					t.Fatalf("expected %s to be in the redirect URIs", uri.URI)
				}
			}
		})
	}
}
