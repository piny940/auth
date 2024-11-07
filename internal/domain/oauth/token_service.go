package oauth

import (
	"auth/internal/domain"
	"crypto/rsa"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type AccessToken struct {
	Value     string
	ExpiresAt time.Time
}
type IDToken struct {
	Value     string
	ExpiresAt time.Time
}

type TokenService struct {
	rsaPrivateKey *rsa.PrivateKey
	issuer        string
	userRepo      domain.IUserRepo
}

const (
	ACCESS_TOKEN_TTL     = 24 * 30 * time.Hour // 30 days
	ACCESS_TOKEN_JTI_LEN = 32
	ID_TOKEN_TTL         = 24 * 3 * time.Hour // 3 days
	ID_TOKEN_JTI_LEN     = 32
)

func NewTokenService(config *Config, userRepo domain.IUserRepo) *TokenService {
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(
		[]byte(config.RsaPrivateKey),
		config.RsaPrivateKeyPassphrase,
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
	jti, err := randomString(ACCESS_TOKEN_JTI_LEN)
	if err != nil {
		return nil, err
	}
	strScopes := make([]string, 0, len(authCode.Scopes))
	for _, s := range authCode.Scopes {
		strScopes = append(strScopes, string(s))
	}
	user, err := s.userRepo.FindByID(authCode.UserID)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(ACCESS_TOKEN_TTL)
	raw := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   s.issuer,
		"exp":   expiresAt.Unix(),
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
		Value:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *TokenService) IssueIDToken(authCode *AuthCode) (*IDToken, error) {
	jti, err := randomString(ID_TOKEN_JTI_LEN)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByID(authCode.UserID)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(ID_TOKEN_TTL)
	raw := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": s.issuer,
		"sub": fmt.Sprintf("id:%d;name:%s", user.ID, user.Name),
		"aud": authCode.ClientID,
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
		"jti": jti,
	})
	token, err := raw.SignedString(s.rsaPrivateKey)
	if err != nil {
		return nil, err
	}
	return &IDToken{
		Value:     token,
		ExpiresAt: expiresAt,
	}, nil
}
