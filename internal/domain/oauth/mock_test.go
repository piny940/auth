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

type approvalRepo struct {
	lastId    ApprovalID
	Approvals []*Approval
}

var _ IApprovalRepo = &approvalRepo{}

func (a *approvalRepo) Approve(clientID ClientID, userID domain.UserID, scopes []TypeScope) error {
	a.Approvals = append(a.Approvals, &Approval{
		ID:        a.lastId,
		ClientID:  clientID,
		UserID:    userID,
		Scopes:    scopes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	a.lastId++
	return nil
}

func (a *approvalRepo) Find(clientID ClientID, userID domain.UserID) (*Approval, error) {
	for _, approval := range a.Approvals {
		if approval.ClientID == clientID && approval.UserID == userID {
			return approval, nil
		}
	}
	return nil, domain.ErrRecordNotFound
}

type clientRepo struct{}

var _ IClientRepo = &clientRepo{}

func (c *clientRepo) Create(client *ClientInput) error {
	panic("unimplemented")
}
func (c *clientRepo) FindWithUserID(id ClientID, userID domain.UserID) (*Client, error) {
	panic("unimplemented")
}
func (c *clientRepo) Update(client *ClientInput, userID domain.UserID) error {
	panic("unimplemented")
}
func (c *clientRepo) FindByID(id ClientID) (*Client, error) {
	panic("unimplemented")
}
func (c *clientRepo) Delete(id ClientID, userID domain.UserID) error {
	panic("unimplemented")
}
func (c *clientRepo) List(userID domain.UserID) ([]*Client, error) {
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
