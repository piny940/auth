package oauth

import (
	"auth/internal/domain"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestNewTokenService(t *testing.T) {
	const issuer = "auth.example.com"
	tokenSvc := NewTokenService(&Config{
		RsaPrivateKey:           RSA_PRIVATE_KEY,
		RsaPrivateKeyPassphrase: RSA_PASS_PHRASE,
		Issuer:                  issuer,
	}, &userRepo{})
	if tokenSvc.issuer != issuer {
		t.Errorf("expected: %s, got: %s", issuer, tokenSvc.issuer)
	}
}

func TestTokenServiceIssueAccessToken(t *testing.T) {
	const TIME_GAP = time.Minute // 許容する誤差
	const issuer = "auth.example.com"
	const userID = 1
	const userName = "test"
	users := []*domain.User{
		{ID: userID, Name: userName, EncryptedPassword: "test"},
	}
	tokenSvc := NewTokenService(&Config{
		RsaPrivateKey:           RSA_PRIVATE_KEY,
		RsaPrivateKeyPassphrase: RSA_PASS_PHRASE,
		RsaKeyId:                RSA_KEY_ID,
		Issuer:                  issuer,
	}, &userRepo{Users: users})
	authCode := &AuthCode{
		UserID: 1,
		Scopes: []TypeScope{ScopeOpenID, ScopeOpenID},
	}
	token, err := tokenSvc.IssueAccessToken(authCode)
	if err != nil {
		t.Fatal(err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(RSA_PUBLIC_KEY))
	if err != nil {
		t.Fatal(err)
	}
	tok, err := jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	claims := tok.Claims.(jwt.MapClaims)
	if claims["iss"] != issuer {
		t.Errorf("expected: %s, got: %s", issuer, claims["iss"])
	}
	if !(time.Now().Add(ACCESS_TOKEN_TTL).Add(-TIME_GAP).Unix() < int64(claims["exp"].(float64)) &&
		int64(claims["exp"].(float64)) < time.Now().Add(ACCESS_TOKEN_TTL).Add(TIME_GAP).Unix()) {
		t.Errorf("exp is invalid")
	}
	if claims["sub"] != fmt.Sprintf("id:%d;name:%s", userID, userName) {
		t.Errorf("sub is invalid. expected: id:%d;name:%s, got: %s", userID, userName, claims["sub"])
	}
	if !(time.Now().Add(-TIME_GAP).Unix() < int64(claims["iat"].(float64)) &&
		int64(claims["iat"].(float64)) < time.Now().Add(TIME_GAP).Unix()) {
		t.Errorf("iat is invalid")
	}
	if len(claims["jti"].(string)) != ACCESS_TOKEN_JTI_LEN {
		t.Errorf("jti is invalid")
	}
	if claims["scope"] != fmt.Sprintf("%s %s", ScopeOpenID, ScopeOpenID) {
		t.Errorf("scope is invalid")
	}
	if tok.Header["kid"] != RSA_KEY_ID {
		t.Errorf("kid is invalid")
	}
}

func TestTokenServiceIssueIDToken(t *testing.T) {
	const TIME_GAP = time.Minute // 許容する誤差
	const issuer = "auth.example.com"
	const userID = 1
	const userName = "test"
	const clientID = "test_client"
	users := []*domain.User{
		{ID: userID, Name: userName, EncryptedPassword: "test"},
	}
	tokenSvc := NewTokenService(&Config{
		RsaPrivateKey:           RSA_PRIVATE_KEY,
		RsaPrivateKeyPassphrase: RSA_PASS_PHRASE,
		RsaKeyId:                RSA_KEY_ID,
		Issuer:                  issuer,
	}, &userRepo{Users: users})
	authCode := &AuthCode{
		UserID:   userID,
		ClientID: clientID,
		Scopes:   []TypeScope{ScopeOpenID, ScopeOpenID},
	}
	token, err := tokenSvc.IssueIDToken(authCode)
	if err != nil {
		t.Fatal(err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(RSA_PUBLIC_KEY))
	if err != nil {
		t.Fatal(err)
	}
	tok, err := jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	claims := tok.Claims.(jwt.MapClaims)
	if claims["iss"] != issuer {
		t.Errorf("expected: %s, got: %s", issuer, claims["iss"])
	}
	if !(time.Now().Add(ID_TOKEN_TTL).Add(-TIME_GAP).Unix() < int64(claims["exp"].(float64)) &&
		int64(claims["exp"].(float64)) < time.Now().Add(ID_TOKEN_TTL).Add(TIME_GAP).Unix()) {
		t.Errorf("exp is invalid")
	}
	if claims["sub"] != fmt.Sprintf("id:%d;name:%s", userID, userName) {
		t.Errorf("sub is invalid. expected: id:%d;name:%s, got: %s", userID, userName, claims["sub"])
	}
	if claims["aud"] != clientID {
		t.Errorf("aud is invalid")
	}
	if !(time.Now().Add(-TIME_GAP).Unix() < int64(claims["iat"].(float64)) &&
		int64(claims["iat"].(float64)) < time.Now().Add(TIME_GAP).Unix()) {
		t.Errorf("iat is invalid")
	}
	if len(claims["jti"].(string)) != ACCESS_TOKEN_JTI_LEN {
		t.Errorf("jti is invalid")
	}
	if tok.Header["kid"] != RSA_KEY_ID {
		t.Errorf("kid is invalid")
	}
}

const (
	RSA_PRIVATE_KEY = `-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: AES-128-CBC,B6F46C744F57D37FB96E49CC625C7BA5

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
-----END RSA PRIVATE KEY-----
`

	RSA_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAqSqLJDQio/mpbWeBnCCm
VrW+3kvreuF26ngwDr2b/PvDsYlgk+7oICgiHAyGrh87glVXuI7QZE/8/A7f9owj
yFF9gS+ev/NdFKk72SfctqFfts+E78nnujnGtariV3YZds6LRnYi1op/mySzrmnV
cHUCUoLpbsYTisCdjfdL4tSvJ5kfs4rebiFCHxoJ2c00KCgb4Y+cPFEy4ncoxD6s
swWPu/rn2tyj7zjYoUqvsYScMyzYxouS3HateN+dBYPYX7zKGUk5PVL4yySfkutv
eEikRU0WmMyje9g33GLnSMgBPtXbkXOxjkN5AS5nvnZ023gLf8urlgQlzuRa8oLI
oGfbmR5J1YkYT0AsNYPy/JAzsFQlc23GQrTv4ay3PU23jmtgcHRF97hdV570ZMBn
AkQbTuNJ7ogkg/NarUvd/nZ9HaqhQikwexjSFoDahorpYW3kNbfA69V+f3JzRtK/
V2PF+Zp4YJwd/dAu0kUnMNHzVMjqO4nG29+B61MTvp/nAgMBAAE=
-----END PUBLIC KEY-----
`
	RSA_PASS_PHRASE = "password"
	RSA_KEY_ID      = "qf7CErL9vRujPENk5"
)
