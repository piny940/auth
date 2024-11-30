package api

import (
	"context"

	"github.com/labstack/echo/v4"
)

const CTX_COOKIE_KEY = "cookie"

func (s *Server) SessionInterfaceLogin(ctx context.Context, request SessionInterfaceLoginRequestObject) (SessionInterfaceLoginResponseObject, error) {
	user, err := s.UserUsecase.Login(request.Body.Name, request.Body.Password)
	if err != nil {
		s.logger.Infof("failed to login: %v", err)
		return SessionInterfaceLogin400JSONResponse{
			Error:            InvalidNameOrPassword,
			ErrorDescription: "name or password is incorrect",
		}, nil
	}
	if err := s.Auth.Login(ctx, user); err != nil {
		s.logger.Errorf("failed to login: %v", err)
		return nil, err
	}
	return &SessionInterfaceLogin204Response{}, nil
}

// SessionInterfaceLogout implements StrictServerInterface.
func (s *Server) SessionInterfaceLogout(ctx context.Context, request SessionInterfaceLogoutRequestObject) (SessionInterfaceLogoutResponseObject, error) {
	session, err := s.Auth.CurrentUser(ctx)
	if err != nil {
		s.logger.Errorf("failed to get current user: %v", err)
		return nil, err
	}
	if session == nil {
		return nil, echo.ErrUnauthorized
	}
	if err := s.Auth.Logout(ctx); err != nil {
		s.logger.Errorf("failed to logout: %v", err)
		return nil, err
	}
	return &SessionInterfaceLogout204Response{}, nil
}

// SessionInterfaceMe implements StrictServerInterface.
func (s *Server) SessionInterfaceMe(ctx context.Context, request SessionInterfaceMeRequestObject) (SessionInterfaceMeResponseObject, error) {
	session, err := s.Auth.CurrentUser(ctx)
	if err != nil {
		s.logger.Errorf("failed to get current user: %v", err)
		return nil, err
	}
	return &SessionInterfaceMe200JSONResponse{
		User: &User{
			Id:   int64(session.User.ID),
			Name: session.User.Name,
		},
	}, nil
}
