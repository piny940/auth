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
