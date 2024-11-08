package gateway

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"errors"
	"slices"
	"time"

	"gorm.io/gorm"
)

type AuthCodeRepo struct {
	db    *infrastructure.DB
	query *query.Query
}

const AUTH_CODE_SCOPE_BATCH_SIZE = 20

func NewAuthCodeRepo(db *infrastructure.DB) oauth.IAuthCodeRepo {
	query := query.Use(db.Client)
	return &AuthCodeRepo{
		db:    db,
		query: query,
	}
}

func (a *AuthCodeRepo) Find(value string) (*oauth.AuthCode, error) {
	authCode, err := a.query.AuthCode.Where(a.query.AuthCode.Value.Eq(value)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	mScopes, err := a.query.AuthCodeScope.Where(
		a.query.AuthCodeScope.AuthCodeID.Eq(authCode.ID),
	).Find()
	if err != nil {
		return nil, err
	}
	return toDomainAuthCode(authCode, mScopes), nil
}

func (a *AuthCodeRepo) Create(value string, clientID oauth.ClientID, userID domain.UserID, scopes []oauth.TypeScope, expiresAt time.Time, redirectURI string) error {
	err := a.query.AuthCode.Create(&model.AuthCode{
		Value:       value,
		ClientID:    string(clientID),
		UserID:      int64(userID),
		ExpiresAt:   expiresAt,
		Used:        false,
		RedirectURI: redirectURI,
	})
	if err != nil {
		return err
	}
	code, err := a.query.AuthCode.Where(a.query.AuthCode.Value.Eq(value)).First()
	if err != nil {
		return err
	}
	compactScopes := make([]oauth.TypeScope, 0, len(scopes))
	for _, s := range scopes {
		if !slices.Contains(compactScopes, s) {
			compactScopes = append(compactScopes, s)
		}
	}
	adds := make([]*model.AuthCodeScope, 0, len(compactScopes))
	for _, s := range compactScopes {
		adds = append(adds, &model.AuthCodeScope{
			ScopeID:    ScopeMapReverse[s],
			AuthCodeID: code.ID,
		})
	}
	if err := a.query.AuthCodeScope.CreateInBatches(adds, AUTH_CODE_SCOPE_BATCH_SIZE); err != nil {
		return err
	}
	return nil
}

func (a *AuthCodeRepo) Use(value string) error {
	_, err := a.query.AuthCode.Where(a.query.AuthCode.Value.Eq(value)).Update(
		a.query.AuthCode.Used, true,
	)
	return err
}

func toDomainAuthCode(m *model.AuthCode, mScopes []*model.AuthCodeScope) *oauth.AuthCode {
	scopes := make([]oauth.TypeScope, 0, len(mScopes))
	for _, s := range mScopes {
		scopes = append(scopes, ScopeMap[s.ScopeID])
	}
	return &oauth.AuthCode{
		Value:       m.Value,
		ClientID:    oauth.ClientID(m.ClientID),
		UserID:      domain.UserID(m.UserID),
		ExpiresAt:   m.ExpiresAt,
		Used:        m.Used,
		RedirectURI: m.RedirectURI,
		Scopes:      scopes,
	}
}
