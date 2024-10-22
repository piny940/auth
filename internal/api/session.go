package api

import (
	"context"
)

const CTX_COOKIE_KEY = "cookie"

func (s *Server) SessionInterfaceLogin(ctx context.Context, request SessionInterfaceLoginRequestObject) (SessionInterfaceLoginResponseObject, error) {
	user, err := s.AuthUsecase.Login(request.Body.Name, request.Body.Password)
	if err != nil {
		return SessionInterfaceLogin400JSONResponse{
			Error:            InvalidNameOrPassword,
			ErrorDescription: "name or password is incorrect",
		}, nil
	}
	cookie, err := Login(ctx, user)
	if err != nil {
		return nil, err
	}
	return &SessionInterfaceLogin204Response{
		Headers: SessionInterfaceLogin204ResponseHeaders{
			SetCookie: cookie.String(),
		},
	}, nil
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
