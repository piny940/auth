package oauth

import (
	"crypto/rsa"
	"testing"
)

func TestIssueJwks(t *testing.T) {
	conf := &Config{
		RsaPublicKey: RSA_PUBLIC_KEY,
		RsaKeyId:     RSA_KEY_ID,
	}
	jwkSvc := NewJWKsService(conf)
	jwks, err := jwkSvc.IssueJwks()
	if err != nil {
		t.Fatal(err)
	}
	jwk, ok := jwks.LookupKeyID(RSA_KEY_ID)
	if !ok {
		t.Errorf("JWK not found: %s", RSA_KEY_ID)
	}
	var raw rsa.PublicKey
	if err := jwk.Raw(&raw); err != nil {
		t.Fatal(err)
	}
}
