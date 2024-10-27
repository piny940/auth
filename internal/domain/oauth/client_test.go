package oauth

import (
	"errors"
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
