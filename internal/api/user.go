package api

import (
	"auth/internal/domain"
	"context"
	"errors"
)

func (s *Server) UsersInterfaceSignup(ctx context.Context, request UsersInterfaceSignupRequestObject) (UsersInterfaceSignupResponseObject, error) {
	_, err := s.UserUsecase.SignUp(
		request.Body.Email, request.Body.Name, request.Body.Password, request.Body.PasswordConfirmation,
	)
	if errors.Is(err, domain.ErrNameLengthNotEnough) {
		s.logger.Infof("name length not enough: %v", err)
		return UsersInterfaceSignup400JSONResponse{
			Error:            NameLengthNotEnough,
			ErrorDescription: err.Error(),
		}, nil
	}
	if errors.Is(err, domain.ErrNameAlreadyUsed) {
		s.logger.Infof("name already used: %v", err)
		return UsersInterfaceSignup400JSONResponse{
			Error:            NameAlreadyUsed,
			ErrorDescription: err.Error(),
		}, nil
	}
	if errors.Is(err, domain.ErrPasswordLengthNotEnough) {
		s.logger.Infof("password length not enough: %v", err)
		return UsersInterfaceSignup400JSONResponse{
			Error:            PasswordLengthNotEnough,
			ErrorDescription: err.Error(),
		}, nil
	}
	if errors.Is(err, domain.ErrPasswordConfirmation) {
		s.logger.Infof("password confirmation not match: %v", err)
		return UsersInterfaceSignup400JSONResponse{
			Error:            PasswordConfirmationNotMatch,
			ErrorDescription: err.Error(),
		}, nil
	}
	if err != nil {
		s.logger.Errorf("failed to signup: %v", err)
		return nil, err
	}
	return &UsersInterfaceSignup204Response{}, nil
}
