package api

import (
	"context"
	"fmt"
	"log/slog"
)

// V1Login implements StrictServerInterface.
func (s *Server) V1Login(ctx context.Context, request V1LoginRequestObject) (V1LoginResponseObject, error) {
	_, err := s.AuthUsecase.Login(request.Body.Name, request.Body.Password)
	if err != nil {
		slog.Info(fmt.Sprintf("login failed: %s", err.Error()))
		return V1Login400JSONResponse{
			Error:            400,
			ErrorDescription: "usename or password is incorrect",
		}, nil
	}

	// TODO: login
	return V1Login200Response{}, nil
}
