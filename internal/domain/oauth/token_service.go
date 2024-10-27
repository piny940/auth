package oauth

import (
	"auth/internal/domain"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
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
	userRepo      domain.IUserRepo
}

const (
	ACCESS_TOKEN_TTL     = 24 * 30 * time.Hour
	ACCESS_TOKEN_JTI_LEN = 32
)

func NewTokenService(config *Config, userRepo domain.IUserRepo) *TokenService {
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
		userRepo:      userRepo,
	}
}

func (s *TokenService) IssueAccessToken(authCode *AuthCode) (*AccessToken, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, ACCESS_TOKEN_JTI_LEN)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	var jti string
	for _, v := range b {
		jti += string(letters[int(v)%len(letters)])
	}
	strScopes := make([]string, 0, len(authCode.Scopes))
	for _, s := range authCode.Scopes {
		strScopes = append(strScopes, string(s))
	}
	user, err := s.userRepo.FindByID(authCode.UserID)
	if err != nil {
		return nil, err
	}
	raw := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   s.issuer,
		"exp":   time.Now().Add(ACCESS_TOKEN_TTL).Unix(),
		"iat":   time.Now().Unix(),
		"sub":   fmt.Sprintf("id:%d;name:%s", user.ID, user.Name),
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
