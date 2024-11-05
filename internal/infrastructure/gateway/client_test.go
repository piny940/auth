package gateway

import (
	"auth/internal/domain"
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
			uriStrs := make([]string, 0, len(uris))
			for _, uri := range uris {
				uriStrs = append(uriStrs, uri.URI)
			}
			if !redirectUrisEqual(t, s.redirectURIs, uriStrs) {
				t.Fatalf("expected %v, got %v", s.redirectURIs, uriStrs)
			}
		})
	}
}

func TestClientFindByID(t *testing.T) {
	const userID = 45328
	const clientID = "client1"
	initialUsers := []*model.User{
		{Name: "test1", ID: userID, EncryptedPassword: "password"},
	}
	suites := []struct {
		name     string
		clients  []*oauth.ClientInput
		expected error
	}{
		{"with a URI", []*oauth.ClientInput{{RedirectURIs: []string{"http://example.com"}, ID: clientID, UserID: userID}}, nil},
		{"with two URIs", []*oauth.ClientInput{{RedirectURIs: []string{"http://example.com", "http://example1.com"}, ID: clientID, UserID: userID}}, nil},
		{"with no URI", []*oauth.ClientInput{{RedirectURIs: []string{}, ID: clientID, UserID: userID}}, nil},
		{"not found", []*oauth.ClientInput{}, domain.ErrRecordNotFound},
	}
	for _, s := range suites {
		t.Run(s.name, func(t *testing.T) {
			setup(t)
			db := infrastructure.GetDB()
			query := query.Use(db.Client)
			if err := query.User.CreateInBatches(initialUsers, len(initialUsers)); err != nil {
				t.Fatal(err)
			}
			svc := NewClientRepo(db)
			for _, client := range s.clients {
				err := svc.Create(client)
				if err != nil {
					t.Fatal(err)
				}
			}
			actual, err := svc.FindByID(clientID)
			if err != s.expected {
				t.Fatalf("expected %v, got %v", s.expected, err)
			}
			if err != nil {
				return
			}
			if !redirectUrisEqual(t, actual.RedirectURIs, s.clients[0].RedirectURIs) {
				t.Fatalf("expected %v, got %v", s.clients[0].RedirectURIs, actual.RedirectURIs)
			}
		})
	}
}

func TestClientFindWithUserID(t *testing.T) {
	const userID = 45328
	const user2ID = 45329
	const clientID = "client1"
	initialUsers := []*model.User{
		{Name: "test1", ID: userID, EncryptedPassword: "password"},
		{Name: "test2", ID: user2ID, EncryptedPassword: "password"},
	}
	suites := []struct {
		name     string
		clients  []*oauth.ClientInput
		expected error
	}{
		{"with a URI", []*oauth.ClientInput{{RedirectURIs: []string{"http://example.com"}, ID: clientID, UserID: userID}}, nil},
		{"with two URIs", []*oauth.ClientInput{{RedirectURIs: []string{"http://example.com", "http://example1.com"}, ID: clientID, UserID: userID}}, nil},
		{"with no URI", []*oauth.ClientInput{{RedirectURIs: []string{}, ID: clientID, UserID: userID}}, nil},
		{"no clients", []*oauth.ClientInput{}, domain.ErrRecordNotFound},
		{"other user's client", []*oauth.ClientInput{{ID: clientID, UserID: user2ID}}, domain.ErrRecordNotFound},
	}
	for _, s := range suites {
		t.Run(s.name, func(t *testing.T) {
			setup(t)
			db := infrastructure.GetDB()
			query := query.Use(db.Client)
			if err := query.User.CreateInBatches(initialUsers, len(initialUsers)); err != nil {
				t.Fatal(err)
			}
			svc := NewClientRepo(db)
			for _, client := range s.clients {
				err := svc.Create(client)
				if err != nil {
					t.Fatal(err)
				}
			}
			actual, err := svc.FindWithUserID(clientID, userID)
			if err != s.expected {
				t.Fatalf("expected %v, got %v", s.expected, err)
			}
			if err != nil {
				return
			}
			if !redirectUrisEqual(t, actual.RedirectURIs, s.clients[0].RedirectURIs) {
				t.Fatalf("expected %v, got %v", s.clients[0].RedirectURIs, actual.RedirectURIs)
			}
		})
	}
}

func TestClientUpdate(t *testing.T) {
	const userID = 45328
	const clientID = "client1"
	const client2ID = "client2"
	initialUsers := []*model.User{
		{Name: "test1", ID: userID, EncryptedPassword: "password"},
	}
	initialClients := []*oauth.ClientInput{
		{ID: clientID, UserID: userID, RedirectURIs: []string{"http://example.com", "http://example1.com"}},
		{ID: client2ID, UserID: userID, RedirectURIs: []string{"http://example1.com"}},
	}
	suites := []struct {
		name        string
		clientInput *oauth.ClientInput
	}{
		{"change name", &oauth.ClientInput{Name: "new name", RedirectURIs: []string{"http://example.com", "http://example1.com"}, ID: clientID, UserID: userID}},
		{"remove a redirect URI", &oauth.ClientInput{RedirectURIs: []string{"http://example.com"}, ID: clientID, UserID: userID}},
		{"remove all redirect URI", &oauth.ClientInput{RedirectURIs: []string{}, ID: clientID, UserID: userID}},
		{"change a redirect URI", &oauth.ClientInput{RedirectURIs: []string{"http://example1.com", "http://example2.com"}, ID: clientID, UserID: userID}},
		{"change all redirect URI", &oauth.ClientInput{RedirectURIs: []string{"http://example2.com", "http://example3.com"}, ID: clientID, UserID: userID}},
		{"change and remove redirect URI", &oauth.ClientInput{RedirectURIs: []string{"http://example2.com"}, ID: clientID, UserID: userID}},
		{"add redirect URI", &oauth.ClientInput{RedirectURIs: []string{"http://example.com", "http://example1.com", "http://example2.com"}, ID: clientID, UserID: userID}},
		{"change name and redirect URI", &oauth.ClientInput{Name: "new name", RedirectURIs: []string{"http://example.com", "http://example2.com"}, ID: clientID, UserID: userID}},
	}
	for _, s := range suites {
		t.Run(s.name, func(t *testing.T) {
			setup(t)
			db := infrastructure.GetDB()
			query := query.Use(db.Client)
			if err := query.User.CreateInBatches(initialUsers, len(initialUsers)); err != nil {
				t.Fatal(err)
			}
			svc := NewClientRepo(db)
			for _, c := range initialClients {
				if err := svc.Create(c); err != nil {
					t.Fatal(err)
				}
			}
			if err := svc.Update(s.clientInput, userID); err != nil {
				t.Fatal(err)
			}
			client, err := svc.FindByID(clientID)
			if err != nil {
				t.Fatal(err)
			}
			if client.Name != s.clientInput.Name {
				t.Fatalf("expected %s, got %s", s.clientInput.Name, client.Name)
			}
			if client.UserID != s.clientInput.UserID {
				t.Fatalf("expected %d, got %d", userID, client.UserID)
			}
			if !redirectUrisEqual(t, client.RedirectURIs, s.clientInput.RedirectURIs) {
				t.Fatalf("expected %v, got %v", s.clientInput.RedirectURIs, client.RedirectURIs)
			}
		})
	}
}

func TestClientDelete(t *testing.T) {
	const userID = 45328
	const clientID = "client1"
	initialUsers := []*model.User{
		{Name: "test1", ID: userID, EncryptedPassword: "password"},
	}
	initialClients := []*oauth.ClientInput{
		{ID: clientID, UserID: userID, RedirectURIs: []string{"http://example.com", "http://example1.com"}},
	}
	suites := []struct {
		name     string
		clientID oauth.ClientID
		expected error
	}{
		{"existing client", clientID, nil},
	}
	for _, s := range suites {
		t.Run(s.name, func(t *testing.T) {
			setup(t)
			db := infrastructure.GetDB()
			query := query.Use(db.Client)
			if err := query.User.CreateInBatches(initialUsers, len(initialUsers)); err != nil {
				t.Fatal(err)
			}
			svc := NewClientRepo(db)
			for _, c := range initialClients {
				if err := svc.Create(c); err != nil {
					t.Fatal(err)
				}
			}
			err := svc.Delete(s.clientID, userID)
			if err != s.expected {
				t.Fatalf("expected %v, got %v", s.expected, err)
			}
			if err != nil {
				return
			}
			_, err = svc.FindByID(s.clientID)
			if err != domain.ErrRecordNotFound {
				t.Fatalf("expected %v, got %v", domain.ErrRecordNotFound, err)
			}
			uris, err := query.RedirectURI.Where(query.RedirectURI.ClientID.Eq(string(s.clientID))).Find()
			if err != nil {
				t.Fatal(err)
			}
			if len(uris) != 0 {
				t.Fatalf("expected 0 redirect URIs, got %d", len(uris))
			}
		})
	}
}

func TestClientList(t *testing.T) {
	const userID = 45328
	const user2ID = 45329
	const user3ID = 45330
	initialUsers := []*model.User{
		{Name: "with clients", ID: userID, EncryptedPassword: "password"},
		{Name: "with no client", ID: user2ID, EncryptedPassword: "password"},
		{Name: "dummy user", ID: user3ID, EncryptedPassword: "password"},
	}
	initialClients := []*oauth.ClientInput{
		{ID: "with two uris", UserID: userID, RedirectURIs: []string{"http://example.com", "http://example1.com"}},
		{ID: "with a uri", UserID: userID, RedirectURIs: []string{"http://example1.com"}},
		{ID: "with no uri", UserID: userID, RedirectURIs: []string{}},
		{ID: "dummy client", UserID: user3ID, RedirectURIs: []string{"http://example1.com"}},
	}
	suites := []struct {
		name        string
		userID      domain.UserID
		expectedIDs []oauth.ClientID
	}{
		{"clients found", userID, []oauth.ClientID{"with two uris", "with a uri", "with no uri"}},
		{"no clients", user2ID, []oauth.ClientID{}},
	}
	for _, s := range suites {
		t.Run(s.name, func(t *testing.T) {
			setup(t)
			db := infrastructure.GetDB()
			query := query.Use(db.Client)
			if err := query.User.CreateInBatches(initialUsers, len(initialUsers)); err != nil {
				t.Fatal(err)
			}
			svc := NewClientRepo(db)
			for _, c := range initialClients {
				if err := svc.Create(c); err != nil {
					t.Fatal(err)
				}
			}
			actual, err := svc.List(s.userID)
			if err != nil {
				t.Fatal(err)
			}
			if len(actual) != len(s.expectedIDs) {
				t.Fatalf("expected %d clients, got %d", len(s.expectedIDs), len(actual))
			}
			for _, c := range actual {
				if !slices.Contains(s.expectedIDs, c.ID) {
					t.Fatalf("expected %s to be in the list", c.ID)
				}
				expected, err := svc.FindByID(c.ID)
				if err != nil {
					t.Fatal(err)
				}
				if !redirectUrisEqual(t, expected.RedirectURIs, c.RedirectURIs) {
					t.Fatalf("expected %v, got %v", expected.RedirectURIs, c.RedirectURIs)
				}
			}
		})
	}

}

func redirectUrisEqual(t *testing.T, a, b []string) bool {
	t.Helper()

	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
