package gateway

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"slices"
	"time"
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

func (a *AuthCodeRepo) Create(value string, clientID oauth.ClientID, userID domain.UserID, scopes []oauth.TypeScope, expiresAt time.Time) error {
	err := a.query.AuthCode.Create(&model.AuthCode{
		Value:     value,
		ClientID:  string(clientID),
		UserID:    int64(userID),
		ExpiresAt: expiresAt,
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
			ScopeID:    scopeMapReverse[s],
			AuthCodeID: code.ID,
		})
	}
	if err := a.query.AuthCodeScope.CreateInBatches(adds, AUTH_CODE_SCOPE_BATCH_SIZE); err != nil {
		return err
	}
	return nil
}
