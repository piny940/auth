package api

import (
	"context"
	"fmt"
)

func (s *Server) UserinfoGetUserinfo(ctx context.Context, request UserinfoGetUserinfoRequestObject) (UserinfoGetUserinfoResponseObject, error) {
	fmt.Println(s.Auth.AccessScopes(ctx))
	return nil, nil
}
