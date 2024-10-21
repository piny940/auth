package infrastructure

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure/query"
)

type ApprovalRepo struct {
	db    *DB
	query *query.Query
}

func NewApprovalRepo(db *DB) oauth.IApprovalRepo {
	query := query.Use(db.Client)
	return &ApprovalRepo{
		db:    db,
		query: query,
	}
}

var _ oauth.IApprovalRepo = &ApprovalRepo{}

// Find implements oauth.IApprovalRepo.
func (a *ApprovalRepo) Find(clientID oauth.ClientID, userID domain.UserID) (*oauth.Approval, error) {
	panic("unimplemented")
}
