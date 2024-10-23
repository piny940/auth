package api

import (
	"auth/internal/usecase"

	"github.com/kelseyhightower/envconfig"
)

type Server struct {
	AuthUsecase   *usecase.AuthUsecase
	ClientUsecase usecase.IClientUsecase
	Conf          *Config
}

type Config struct {
	ServerUrl  string `envconfig:"SERVER_URL" required:"true"`
	LoginUrl   string `split_words:"true" required:"true"`
	ApproveUrl string `split_words:"true" required:"true"`
}

var _ StrictServerInterface = &Server{}

func NewServer(authUsecase *usecase.AuthUsecase, clientUC usecase.IClientUsecase) *Server {
	conf := &Config{}
	err := envconfig.Process("api", conf)
	if err != nil {
		panic(err)
	}
	return &Server{
		AuthUsecase:   authUsecase,
		ClientUsecase: clientUC,
		Conf:          conf,
	}
}
