package api

import (
	"auth/internal/domain/oauth"
	"context"
	"fmt"
	"slices"
)

func (s *Server) UserinfoGetUserinfo(ctx context.Context, request UserinfoGetUserinfoRequestObject) (UserinfoGetUserinfoResponseObject, error) {
	scopes, err := s.Auth.AccessScopes(ctx)
	fmt.Println(scopes)
	if err != nil {
		s.logger.Errorf("failed to get access scopes: %v", err)
		return nil, err
	}
	session, err := s.Auth.CurrentUser(ctx)
	if err != nil {
		s.logger.Errorf("failed to get current user: %v", err)
		return nil, err
	}
	res := UserinfoGetUserinfo200JSONResponse{
		Sub: fmt.Sprintf("id:%d;name:%s", session.User.ID, session.User.Name),
	}
	if slices.Contains(scopes, oauth.ScopeProfile) {
		res.Name = ptr(session.User.Name)
	}
	return res, nil
}
