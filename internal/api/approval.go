package api

import (
	"auth/internal/domain/oauth"
	"context"
	"errors"
	"strings"
)

func (s *Server) ApprovalsInterfaceApprove(ctx context.Context, request ApprovalsInterfaceApproveRequestObject) (ApprovalsInterfaceApproveResponseObject, error) {
	user, err := CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	scopes := make([]oauth.TypeScope, 0)
	for _, s := range strings.Split(request.Body.Scope, " ") {
		scopes = append(scopes, oauth.TypeScope(s))
	}
	err = s.AuthUsecase.Approve(user, oauth.ClientID(request.Body.ClientId), scopes)
	if errors.Is(err, oauth.ErrInvalidClientID) {
		return ApprovalsInterfaceApprove400JSONResponse{
			Error:            ApprovalsApproveErrInvalidClient,
			ErrorDescription: err.Error(),
		}, nil
	}
	if errors.Is(err, oauth.ErrInvalidScope) {
		return ApprovalsInterfaceApprove400JSONResponse{
			Error:            ApprovalsApproveErrInvalidScope,
			ErrorDescription: err.Error(),
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return ApprovalsInterfaceApprove204Response{}, nil
}
