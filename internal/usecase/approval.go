package usecase

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"context"
	"fmt"
)

type ApprovalUsecase struct {
	ApprovalRepo oauth.IApprovalRepo
}

func NewApprovalUsecase(approvalRepo oauth.IApprovalRepo) *ApprovalUsecase {
	return &ApprovalUsecase{
		ApprovalRepo: approvalRepo,
	}
}

func (u *ApprovalUsecase) List(ctx context.Context, userID domain.UserID) ([]*oauth.Approval, error) {
	approvals, err := u.ApprovalRepo.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list approvals: %w", err)
	}
	return approvals, nil
}

func (u *ApprovalUsecase) Delete(ctx context.Context, userID domain.UserID, clientID oauth.ClientID) error {
	if err := u.ApprovalRepo.Delete(ctx, userID, clientID); err != nil {
		return fmt.Errorf("failed to delete approval: %w", err)
	}
	return nil
}
