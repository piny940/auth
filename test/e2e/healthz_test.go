package e2e

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHealthz(t *testing.T) {
	s := newServer(t)
	defer s.Close()

	resp, err := http.Get(fmt.Sprintf("%s/healthz", s.URL))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
