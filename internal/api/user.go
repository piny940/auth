package api

import (
	"auth/internal/domain"
	"context"
	"errors"
)

func (s *Server) UsersInterfaceSignup(ctx context.Context, request UsersInterfaceSignupRequestObject) (UsersInterfaceSignupResponseObject, error) {
	user, err := s.AuthUsecase.SignUp(
		request.Body.Name, request.Body.Password, request.Body.PasswordConfirmation,
	)
	if errors.Is(err, domain.ErrNameLengthNotEnough) {
		return UsersInterfaceSignup400JSONResponse{
			Error:            NameLengthNotEnough,
			ErrorDescription: err.Error(),
		}, nil
	}
	if errors.Is(err, domain.ErrNameAlreadyUsed) {
		return UsersInterfaceSignup400JSONResponse{
			Error:            NameAlreadyUsed,
			ErrorDescription: err.Error(),
		}, nil
	}
	if errors.Is(err, domain.ErrPasswordLengthNotEnough) {
		return UsersInterfaceSignup400JSONResponse{
			Error:            PasswordLengthNotEnough,
			ErrorDescription: err.Error(),
		}, nil
	}
	if errors.Is(err, domain.ErrPasswordConfirmation) {
		return UsersInterfaceSignup400JSONResponse{
			Error:            PasswordConfirmationNotMatch,
			ErrorDescription: err.Error(),
		}, nil
	}
	if err != nil {
		return nil, err
	}
	cookie, err := Login(ctx, user)
	if err != nil {
		return nil, err
	}
	return &UsersInterfaceSignup204Response{
		Headers: UsersInterfaceSignup204ResponseHeaders{
			SetCookie: cookie.String(),
		},
	}, nil
}
