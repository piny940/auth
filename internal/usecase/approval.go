package usecase

import (
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"context"
	"fmt"
)

type ApprovalUsecase struct {
	ApprovalRepo oauth.IApprovalRepo
	ClientRepo   oauth.IClientRepo
}

func NewApprovalUsecase(approvalRepo oauth.IApprovalRepo, clientRepo oauth.IClientRepo) *ApprovalUsecase {
	return &ApprovalUsecase{
		ApprovalRepo: approvalRepo,
		ClientRepo:   clientRepo,
	}
}

func (u *ApprovalUsecase) List(ctx context.Context, userID domain.UserID) ([]*oauth.Approval, []*oauth.Client, error) {
	approvals, err := u.ApprovalRepo.List(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list approvals: %w", err)
	}
	clientIDs := make([]oauth.ClientID, 0, len(approvals))
	for _, approval := range approvals {
		clientIDs = append(clientIDs, approval.ClientID)
	}
	clients, err := u.ClientRepo.ListByIds(ctx, clientIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list clients: %w", err)
	}
	return approvals, clients, nil
}

func (u *ApprovalUsecase) Delete(ctx context.Context, ID oauth.ApprovalID, userID domain.UserID) error {
	if err := u.ApprovalRepo.Delete(ctx, ID, userID); err != nil {
		return fmt.Errorf("failed to delete approval: %w", err)
	}
	return nil
}
