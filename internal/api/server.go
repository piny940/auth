package api

import (
	"auth/internal/usecase"
	"context"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
)

type Server struct {
	UserUsecase   *usecase.UserUsecase
	OAuthUsecase  *usecase.OAuthUsecase
	ClientUsecase usecase.IClientUsecase
	Conf          *Config
	Auth          Auth
	logger        echo.Logger
}

type Config struct {
	ServerUrl  string `envconfig:"SERVER_URL" required:"true"`
	LoginUrl   string `split_words:"true" required:"true"`
	ApproveUrl string `split_words:"true" required:"true"`
}

var _ StrictServerInterface = &Server{}

func NewServer(
	userUsecase *usecase.UserUsecase,
	oauthUsecase *usecase.OAuthUsecase,
	clientUC usecase.IClientUsecase,
	auth Auth,
) *Server {
	conf := &Config{}
	err := envconfig.Process("api", conf)
	if err != nil {
		panic(err)
	}
	return &Server{
		UserUsecase:   userUsecase,
		OAuthUsecase:  oauthUsecase,
		ClientUsecase: clientUC,
		Conf:          conf,
		Auth:          auth,
	}
}

func (s *Server) SetLogger(logger echo.Logger) {
	s.logger = logger
}

func (s *Server) HealthzCheck(ctx context.Context, request HealthzCheckRequestObject) (HealthzCheckResponseObject, error) {
	return HealthzCheck200Response{}, nil
}
