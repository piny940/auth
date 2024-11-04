package e2e

import (
	"testing"
)

func TestSignupLogin(t *testing.T) {
	s := newServer(t)
	defer s.Close()

	login(t, s)
}
