package oauth

import (
	"errors"
	"slices"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestEncryptClientSecret(t *testing.T) {
	suites := []struct {
		Name      string
		Secret    string
		ToCompare string
		Correct   bool
	}{
		{"correct secret", "test_secret", "test_secret", true},
		{"invalid secret", "test_secret, test_secret", "test_secret", false},
	}
	for _, suit := range suites {
		t.Run(suit.Name, func(t *testing.T) {
			hash, err := EncryptClientSecret(suit.Secret)
			if err != nil {
				t.Fatal(err)
			}
			err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(suit.ToCompare))
			if err != nil && !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				t.Fatal(err)
			}
			actualIncorrect := errors.Is(err, bcrypt.ErrMismatchedHashAndPassword)
			if actualIncorrect != !suit.Correct {
				t.Errorf("unexpected result: %v", actualIncorrect)
			}
		})
	}
}

func TestClientCorrectSecret(t *testing.T) {
	secret := "test_secret"
	suites := []struct {
		Name      string
		ToCompare string
	}{
		{"correct secret", secret},
		{"invalid secret", "invalid_secret"},
	}
	for _, suit := range suites {
		t.Run(suit.Name, func(t *testing.T) {
			hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
			if err != nil {
				t.Fatal(err)
			}
			client := &Client{EncryptedSecret: string(hash)}
			err = client.SecretCorrect(suit.ToCompare)
			if suit.ToCompare == secret {
				if err != nil {
					t.Fatal(err)
				}
			} else {
				if !errors.Is(err, ErrInvalidClientSecret) {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestClientRedirectURIValid(t *testing.T) {
	client := &Client{RedirectURIs: []string{"https://example.com/callback"}}
	suites := []struct {
		Name     string
		URI      string
		Expected bool
	}{
		{"valid uri", "https://example.com/callback", true},
		{"query can be added", "https://example.com/callback?query=1", true},
		{"schema not match", "http://example.com/callback", false},
		{"host not match", "https://example.org/callback", false},
		{"path not match", "https://example.com/callback2", false},
		{"port not match", "https://example.com:8080/callback", false},
		{"child path not match", "https://example.com/callback/child", false},
		{"parent path not match", "https://example.com", false},
		{"fragment not match", "https://example.com/callback#fragment", false},
		{"invalid uri", "invalid_uri", false},
		{"relative uri", "/callback", false},
		{"empty uri", "", false},
	}
	for _, suit := range suites {
		t.Run(suit.Name, func(t *testing.T) {
			actual := client.RedirectURIValid(suit.URI)
			if actual != suit.Expected {
				t.Errorf("unexpected result: %v", actual)
			}
		})
	}
}

func TestIssueClientID(t *testing.T) {
	ids := make([]ClientID, 0, 10)
	for i := 0; i < 10; i++ {
		id, err := IssueClientID()
		if err != nil {
			t.Fatal(err)
		}
		if len(id) != CLIENT_ID_LEN {
			t.Errorf("unexpected id length: %d", len(id))
		}
		ids = append(ids, id)
	}
	compact := slices.Compact(ids)
	if len(ids) != len(compact) {
		t.Error("duplicated id")
	}
}
func TestIssueClientSecret(t *testing.T) {
	const n = 10
	secrets := make([]string, 0, n)
	for i := 0; i < n; i++ {
		secret, hash, err := IssueClientSecret()
		if err != nil {
			t.Errorf("failed to issue client secret: %v", err)
		}
		if len(secret) != CLIENT_SECRET_LEN {
			t.Errorf("unexpected secret length: %d", len(secret))
		}
		secrets = append(secrets, secret)
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
		if err != nil {
			t.Errorf("failed to compare hash and secret: %v", err)
		}
	}
	compact := slices.Compact(secrets)
	if len(secrets) != len(compact) {
		t.Error("duplicated secret")
	}
}
