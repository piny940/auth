package oauth

import (
	"auth/internal/domain"
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
	Create(clientID ClientID, userID domain.UserID, scopes []TypeScope) error
}

type ApprovalService struct {
	ApprovalRepo IApprovalRepo
	ClientRepo   IClientRepo
}

func NewApprovalService(approvalRepo IApprovalRepo, clientRepo IClientRepo) *ApprovalService {
	return &ApprovalService{
		ApprovalRepo: approvalRepo,
		ClientRepo:   clientRepo,
	}
}
func (s *ApprovalService) Approved(r *AuthRequest, user *domain.User) (bool, error) {
	approval, err := s.ApprovalRepo.Find(r.ClientID, user.ID)
	if errors.Is(err, domain.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if approval == nil {
		return false, nil
	}
	for _, scope := range r.Scopes {
		if !slices.Contains(approval.Scopes, scope) {
			return false, nil
		}
	}
	return true, nil
}

func (s *ApprovalService) Approve(clientID ClientID, userID domain.UserID, scopes []TypeScope) error {
	_, err := s.ClientRepo.FindByID(clientID)
	if err != nil {
		return err
	}
	if err := ValidScopes(scopes); err != nil {
		return err
	}
	if err := s.ApprovalRepo.Create(clientID, userID, scopes); err != nil {
		return err
	}
	return nil
}
