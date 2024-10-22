package api

import (
	"context"
)

const CTX_COOKIE_KEY = "cookie"

// SessionInterfaceLogin implements StrictServerInterface.
func (s *Server) SessionInterfaceLogin(ctx context.Context, request SessionInterfaceLoginRequestObject) (SessionInterfaceLoginResponseObject, error) {
	panic("unimplemented")
}

// SessionInterfaceLogout implements StrictServerInterface.
func (s *Server) SessionInterfaceLogout(ctx context.Context, request SessionInterfaceLogoutRequestObject) (SessionInterfaceLogoutResponseObject, error) {
	panic("unimplemented")
}

// SessionInterfaceMe implements StrictServerInterface.
func (s *Server) SessionInterfaceMe(ctx context.Context, request SessionInterfaceMeRequestObject) (SessionInterfaceMeResponseObject, error) {
	user, err := CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return &SessionInterfaceMe200JSONResponse{
		User: &User{
			Id:   int64(user.ID),
			Name: user.Name,
		},
	}, nil
}
