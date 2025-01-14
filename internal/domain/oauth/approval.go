package oauth

import (
	"auth/internal/domain"
	"context"
	"errors"
	"slices"
	"time"
)

type ApprovalID int64
type Approval struct {
	ID        ApprovalID
	ClientID  ClientID
	UserID    domain.UserID
	Scopes    []TypeScope
	CreatedAt time.Time
	UpdatedAt time.Time
}
type IApprovalRepo interface {
	Find(clientID ClientID, userID domain.UserID) (*Approval, error)
	Approve(clientID ClientID, userID domain.UserID, scopes []TypeScope) error
	List(ctx context.Context, userID domain.UserID) ([]*Approval, error)
	Delete(ctx context.Context, ID ApprovalID, userID domain.UserID) error
}

type ApprovalService struct {
	ApprovalRepo IApprovalRepo
}

func NewApprovalService(approvalRepo IApprovalRepo) *ApprovalService {
	return &ApprovalService{
		ApprovalRepo: approvalRepo,
	}
}
func (s *ApprovalService) Approved(clientID ClientID, byUserID domain.UserID, scopes []TypeScope) (bool, error) {
	approval, err := s.ApprovalRepo.Find(clientID, byUserID)
	if errors.Is(err, domain.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if approval == nil {
		return false, nil
	}
	for _, scope := range scopes {
		if !slices.Contains(approval.Scopes, scope) {
			return false, nil
		}
	}
	return true, nil
}
