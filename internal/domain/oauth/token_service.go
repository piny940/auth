package oauth

import (
	"crypto/rand"
	"crypto/rsa"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type AccessToken struct {
	Value string
}

type TokenService struct {
	rsaPrivateKey *rsa.PrivateKey
	issuer        string
}

const (
	ACCESS_TOKEN_TTL     = 24 * 30 * time.Hour
	ACCESS_TOKEN_JTI_LEN = 32
)

func NewTokenService() *TokenService {
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(
		[]byte(config.RsaPrivateKey),
		config.RsaPrivateKeyPassPhrase,
	)
	if err != nil {
		panic(err)
	}
	return &TokenService{
		rsaPrivateKey: rsaPrivateKey,
		issuer:        config.Issuer,
	}
}

func (s *TokenService) IssueAccessToken(authCode AuthCode, scopes []TypeScope) (*AccessToken, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, ACCESS_TOKEN_JTI_LEN)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	var jti string
	for _, v := range b {
		jti += string(letters[int(v)%len(letters)])
	}
	strScopes := make([]string, 0, len(scopes))
	for _, s := range scopes {
		strScopes = append(strScopes, string(s))
	}
	raw := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   s.issuer,
		"exp":   time.Now().Add(ACCESS_TOKEN_TTL).Unix(),
		"iat":   time.Now().Unix(),
		"jti":   jti,
		"scope": strings.Join(strScopes, " "),
	})
	token, err := raw.SignedString(s.rsaPrivateKey)
	if err != nil {
		return nil, err
	}

	return &AccessToken{
		Value: token,
	}, nil
}
