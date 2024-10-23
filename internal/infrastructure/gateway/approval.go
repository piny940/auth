package gateway

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"errors"
	"slices"

	"gorm.io/gorm"
)

type ApprovalRepo struct {
	db    *infrastructure.DB
	query *query.Query
}

const APPROVAL_SCOPE_BATCH_SIZE = 20

var scopeMap = map[int32]oauth.TypeScope{
	0: oauth.ScopeOpenID,
}
var scopeMapReverse = map[oauth.TypeScope]int32{
	oauth.ScopeOpenID: 0,
}

func NewApprovalRepo(db *infrastructure.DB) oauth.IApprovalRepo {
	query := query.Use(db.Client)
	return &ApprovalRepo{
		db:    db,
		query: query,
	}
}

var _ oauth.IApprovalRepo = &ApprovalRepo{}

func (a *ApprovalRepo) Find(clientID oauth.ClientID, userID domain.UserID) (*oauth.Approval, error) {
	approval, err := a.query.Approval.Where(
		a.query.Approval.ClientID.Eq(string(clientID)),
		a.query.Approval.UserID.Eq(int64(userID)),
	).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	scopes, err := a.query.ApprovalScope.Where(a.query.ApprovalScope.ApprovalID.Eq(approval.ID)).Find()
	if err != nil {
		return nil, err
	}
	return toDomainApproval(approval, scopes), nil
}

func (a *ApprovalRepo) Create(clientID oauth.ClientID, userID domain.UserID, scopes []oauth.TypeScope) error {
	err := a.query.Approval.Create(&model.Approval{
		ClientID: string(clientID),
		UserID:   int64(userID),
	})
	if err != nil {
		return err
	}
	approval, err := a.query.Approval.Where(
		a.query.Approval.ClientID.Eq(string(clientID)),
		a.query.Approval.UserID.Eq(int64(userID)),
	).First()
	if err != nil {
		return err
	}
	compactScopes := slices.Compact(scopes)
	mScopes := make([]*model.ApprovalScope, 0, len(compactScopes))
	for _, s := range compactScopes {
		mScopes = append(mScopes, &model.ApprovalScope{
			ScopeID:    scopeMapReverse[s],
			ApprovalID: approval.ID,
		})
	}
	if err := a.query.ApprovalScope.CreateInBatches(mScopes, APPROVAL_SCOPE_BATCH_SIZE); err != nil {
		return err
	}
	return nil
}

func toDomainApproval(approval *model.Approval, approvalScopes []*model.ApprovalScope) *oauth.Approval {
	scopes := make([]oauth.TypeScope, 0, len(approvalScopes))
	for _, s := range approvalScopes {
		scopes = append(scopes, scopeMap[s.ScopeID])
	}
	return &oauth.Approval{
		ID:        oauth.ApprovalID(approval.ID),
		ClientID:  oauth.ClientID(approval.ClientID),
		UserID:    domain.UserID(approval.UserID),
		Scopes:    scopes,
		CreatedAt: approval.CreatedAt,
		UpdatedAt: approval.UpdatedAt,
	}
}
