package api

import (
	"auth/internal/domain/oauth"
	"context"
	"errors"
	"strings"
)

func (s *Server) ApprovalsInterfaceApprove(ctx context.Context, request ApprovalsInterfaceApproveRequestObject) (ApprovalsInterfaceApproveResponseObject, error) {
	session, err := s.Auth.CurrentUser(ctx)
	if err != nil {
		s.logger.Errorf("failed to get current user: %v", err)
		return nil, err
	}
	scopes := make([]oauth.TypeScope, 0)
	for _, s := range strings.Split(request.Body.Scope, " ") {
		scopes = append(scopes, oauth.TypeScope(s))
	}
	err = s.OAuthUsecase.Approve(session.User, oauth.ClientID(request.Body.ClientId), scopes)
	if errors.Is(err, oauth.ErrInvalidClientID) {
		s.logger.Infof("invalid client id: %v", err)
		return ApprovalsInterfaceApprove400JSONResponse{
			Error:            ApprovalsApproveErrInvalidClient,
			ErrorDescription: err.Error(),
		}, nil
	}
	if errors.Is(err, oauth.ErrInvalidScope) {
		s.logger.Infof("invalid scope: %v", err)
		return ApprovalsInterfaceApprove400JSONResponse{
			Error:            ApprovalsApproveErrInvalidScope,
			ErrorDescription: err.Error(),
		}, nil
	}
	if err != nil {
		s.logger.Errorf("failed to approve: %v", err)
		return nil, err
	}
	return ApprovalsInterfaceApprove204Response{}, nil
}

// ----------------------------------- dashboard api -------------------------------

func (s *Server) ApprovalsDeleteApproval(ctx context.Context, request ApprovalsDeleteApprovalRequestObject) (ApprovalsDeleteApprovalResponseObject, error) {
	session, err := s.Auth.CurrentUser(ctx)
	if err != nil {
		s.logger.Errorf("failed to get current user: %v", err)
		return nil, err
	}
	if err := s.ApprovalUsecase.Delete(ctx, oauth.ApprovalID(request.Id), session.User.ID); err != nil {
		s.logger.Errorf("failed to delete approval: %v", err)
		return nil, err
	}
	return ApprovalsDeleteApproval204Response{}, nil
}

func (s *Server) ApprovalsListApprovals(ctx context.Context, request ApprovalsListApprovalsRequestObject) (ApprovalsListApprovalsResponseObject, error) {
	session, err := s.Auth.CurrentUser(ctx)
	if err != nil {
		s.logger.Errorf("failed to get current user: %v", err)
		return nil, err
	}
	dApprovals, dClients, err := s.ApprovalUsecase.List(ctx, session.User.ID)
	if err != nil {
		s.logger.Errorf("failed to list approvals: %v", err)
		return nil, err
	}
	approvalByClientID := make(map[oauth.ClientID]*oauth.Approval, len(dClients))
	for _, a := range dApprovals {
		approvalByClientID[a.ClientID] = a
	}
	clientByApprovalID := make(map[oauth.ApprovalID]*oauth.Client, len(dApprovals))
	for _, c := range dClients {
		clientByApprovalID[approvalByClientID[c.ID].ID] = c
	}

	approvals := make([]Approval, 0, len(dApprovals))
	for _, a := range dApprovals {
		scopes := make([]string, 0, len(a.Scopes))
		for _, s := range a.Scopes {
			scopes = append(scopes, string(s))
		}
		client := clientByApprovalID[a.ID]
		approvals = append(approvals, Approval{
			Id:     int64(a.ID),
			Scopes: scopes,
			Client: PublicClient{
				Id:   string(client.ID),
				Name: client.Name,
			},
		})
	}
	return ApprovalsListApprovals200JSONResponse{
		Approvals: approvals,
	}, nil
}
