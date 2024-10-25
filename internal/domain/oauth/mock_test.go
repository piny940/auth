package oauth

import (
	"auth/internal/domain"
	"time"
)

type authCodeRepo struct {
	authCodes []*AuthCode
}

var _ IAuthCodeRepo = &authCodeRepo{}

func (a *authCodeRepo) Find(value string) (*AuthCode, error) {
	for _, code := range a.authCodes {
		if code.Value == value {
			return code, nil
		}
	}
	return nil, domain.ErrRecordNotFound
}

func (a *authCodeRepo) Create(value string, clientID ClientID, userID domain.UserID, scopes []TypeScope, expiresAt time.Time, redirectURI string) error {
	a.authCodes = append(a.authCodes, &AuthCode{
		Value:       value,
		ClientID:    clientID,
		UserID:      userID,
		ExpiresAt:   expiresAt,
		Used:        false,
		RedirectURI: redirectURI,
		Scopes:      scopes,
	})
	return nil
}

type approvalRepo struct{}

var _ IApprovalRepo = &approvalRepo{}

func (a *approvalRepo) Create(clientID ClientID, userID domain.UserID, scopes []TypeScope) error {
	panic("unimplemented")
}

func (a *approvalRepo) Find(clientID ClientID, userID domain.UserID) (*Approval, error) {
	panic("unimplemented")
}

type clientRepo struct{}

var _ IClientRepo = &clientRepo{}

func (c *clientRepo) FindByID(id ClientID) (*Client, error) {
	panic("unimplemented")
}

type userRepo struct{ Users []*domain.User }

var _ domain.IUserRepo = &userRepo{}

// Create implements domain.IUserRepo.
func (u *userRepo) Create(name string, encryptedPassword string) error {
	panic("unimplemented")
}

// FindByID implements domain.IUserRepo.
func (u *userRepo) FindByID(id domain.UserID) (*domain.User, error) {
	for _, user := range u.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, domain.ErrRecordNotFound
}

// FindByName implements domain.IUserRepo.
func (u *userRepo) FindByName(name string) (*domain.User, error) {
	panic("unimplemented")
}
