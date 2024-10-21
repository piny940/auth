package api

import (
	"auth/internal/usecase"

	"github.com/labstack/echo/v4"
)

type Server struct {
	AuthUsecase *usecase.AuthUsecase
}

var _ ServerInterface = &Server{}

func NewServer(authUsecase *usecase.AuthUsecase) *Server {
	return &Server{
		AuthUsecase: authUsecase,
	}
}

// ClientsCreateClient implements ServerInterface.
func (s *Server) ClientsCreateClient(ctx echo.Context) error {
	panic("unimplemented")
}

// ClientsDeleteClient implements ServerInterface.
func (s *Server) ClientsDeleteClient(ctx echo.Context, id int64) error {
	panic("unimplemented")
}

// ClientsListClients implements ServerInterface.
func (s *Server) ClientsListClients(ctx echo.Context, params ClientsListClientsParams) error {
	panic("unimplemented")
}

// ClientsUpdateClient implements ServerInterface.
func (s *Server) ClientsUpdateClient(ctx echo.Context, id int64) error {
	panic("unimplemented")
}

// TokenGetToken implements ServerInterface.
func (s *Server) TokenGetToken(ctx echo.Context) error {
	panic("unimplemented")
}
