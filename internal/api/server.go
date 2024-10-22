package api

import (
	"auth/internal/usecase"

	"github.com/kelseyhightower/envconfig"
)

type Server struct {
	AuthUsecase *usecase.AuthUsecase
	Conf        *Config
}

type Config struct {
	ServerUrl  string `envconfig:"SERVER_URL" required:"true"`
	LoginUrl   string `split_words:"true" required:"true"`
	ApproveUrl string `split_words:"true" required:"true"`
}

var _ StrictServerInterface = &Server{}

func NewServer(authUsecase *usecase.AuthUsecase) *Server {
	conf := &Config{}
	err := envconfig.Process("api", conf)
	if err != nil {
		panic(err)
	}
	return &Server{
		AuthUsecase: authUsecase,
		Conf:        conf,
	}
}

// // PostAuthorize implements ServerInterface.
// func (s *Server) PostAuthorize(ctx echo.Context) error {
// 	panic("unimplemented")
// }
