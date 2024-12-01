package gateway

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"context"
	"errors"
	"fmt"
	"slices"

	"gorm.io/gorm"
)

type ApprovalRepo struct {
	db    *infrastructure.DB
	query *query.Query
}

const APPROVAL_SCOPE_BATCH_SIZE = 20

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

func (a *ApprovalRepo) Approve(clientID oauth.ClientID, userID domain.UserID, scopes []oauth.TypeScope) error {
	var approval *model.Approval
	var err error
	{ // create approval if not exists
		_, err = a.query.Approval.Where(
			a.query.Approval.ClientID.Eq(string(clientID)),
			a.query.Approval.UserID.Eq(int64(userID)),
		).First()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = a.query.Approval.Create(&model.Approval{
				ClientID: string(clientID),
				UserID:   int64(userID),
			})
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}
	approval, err = a.query.Approval.Where(
		a.query.Approval.ClientID.Eq(string(clientID)),
		a.query.Approval.UserID.Eq(int64(userID)),
	).First()
	if err != nil {
		return err
	}
	existMScopes, err := a.query.ApprovalScope.Where(a.query.ApprovalScope.ApprovalID.Eq(approval.ID)).Find()
	if err != nil {
		return err
	}
	existScopes := make([]oauth.TypeScope, 0, len(existMScopes))
	for _, s := range existMScopes {
		existScopes = append(existScopes, ScopeMap[s.ScopeID])
	}
	compactScopes := slices.Compact(scopes)
	adds := make([]oauth.TypeScope, 0)
	for _, s := range compactScopes {
		if !slices.Contains(existScopes, s) {
			adds = append(adds, s)
		}
	}
	mAdds := make([]*model.ApprovalScope, 0, len(adds))
	for _, s := range adds {
		scopeID, ok := ScopeMapReverse[s]
		if !ok {
			return oauth.ErrInvalidScope
		}
		mAdds = append(mAdds, &model.ApprovalScope{
			ScopeID:    scopeID,
			ApprovalID: approval.ID,
		})
	}
	if err := a.query.ApprovalScope.CreateInBatches(mAdds, APPROVAL_SCOPE_BATCH_SIZE); err != nil {
		return err
	}
	return nil
}

func (a *ApprovalRepo) List(ctx context.Context, userID domain.UserID) ([]*oauth.Approval, error) {
	approvals, err := a.query.Approval.Where(
		a.query.Approval.UserID.Eq(int64(userID)),
	).WithContext(ctx).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to list approvals: %w", err)
	}
	approvalIds := make([]int64, 0, len(approvals))
	for _, a := range approvals {
		approvalIds = append(approvalIds, a.ID)
	}
	approvalScopes, err := a.query.ApprovalScope.Where(
		a.query.ApprovalScope.ApprovalID.In(approvalIds...),
	).WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}
	scopesMap := make(map[int64][]*model.ApprovalScope)
	for _, s := range approvalScopes {
		if _, ok := scopesMap[s.ApprovalID]; !ok {
			scopesMap[s.ApprovalID] = make([]*model.ApprovalScope, 0)
		}
		scopesMap[s.ApprovalID] = append(scopesMap[s.ApprovalID], s)
	}
	dApprovals := make([]*oauth.Approval, 0, len(approvals))
	for _, a := range approvals {
		dApprovals = append(dApprovals, toDomainApproval(a, scopesMap[a.ID]))
	}
	return dApprovals, nil
}

func (a *ApprovalRepo) Delete(ctx context.Context, userID domain.UserID, clientID oauth.ClientID) error {
	panic("unimplemented")
}

func toDomainApproval(approval *model.Approval, approvalScopes []*model.ApprovalScope) *oauth.Approval {
	scopes := make([]oauth.TypeScope, 0, len(approvalScopes))
	for _, s := range approvalScopes {
		scopes = append(scopes, ScopeMap[s.ScopeID])
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
