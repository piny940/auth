package oauth

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
)

type JWKsService struct {
	rsaPublicKey *rsa.PublicKey
	rsaKeyId     string
}

func NewJWKsService(conf *Config) *JWKsService {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(conf.RsaPublicKey))
	if err != nil {
		panic(err)
	}
	return &JWKsService{
		rsaPublicKey: key,
		rsaKeyId:     conf.RsaKeyId,
	}
}

func (s *JWKsService) IssueJwks() (jwk.Set, error) {
	key, err := jwk.New(s.rsaPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWKs: %w", err)
	}
	if _, ok := key.(jwk.RSAPrivateKey); ok {
		return nil, fmt.Errorf("failed to create JWKs: %s", ErrInvalidKeyType)
	}
	key.Set(jwk.KeyIDKey, s.rsaKeyId)
	jwks := jwk.NewSet()
	jwks.Add(key)
	return jwks, nil
}

var ErrInvalidKeyType = errors.New("invalid key type")
