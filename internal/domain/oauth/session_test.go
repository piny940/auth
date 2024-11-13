package oauth

import (
	"errors"
	"testing"
	"time"
)

func TestSessionValid(t *testing.T) {
	suites := []struct {
		name     string
		authTime time.Time
		maxAge   time.Duration
		expect   error
	}{
		{"expired(maxAge: 0)", time.Now(), 0, ErrSessionExpired},
		{"expired(maxAge: 1s)", time.Now().Add(-1 * time.Second), 1 * time.Second, ErrSessionExpired},
		{"ok", time.Now().Add(-1 * time.Second), 2 * time.Second, nil},
	}
	for _, suit := range suites {
		t.Run(suit.name, func(t *testing.T) {
			s := &Session{AuthTime: suit.authTime}
			err := s.Valid(suit.maxAge)
			if !errors.Is(err, suit.expect) {
				t.Errorf("expected: %v, got: %v", suit.expect, err)
			}
		})
	}
}
