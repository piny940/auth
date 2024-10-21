//go:build wireinject

package di

import (
	"auth/internal/api"
	"auth/internal/infrastructure"
	"auth/internal/usecase"

	"github.com/google/wire"
)

func (r *registry) NewServer() *api.Server {
	wire.Build(
		api.NewServer,
		usecase.NewAuthUsecase,
		infrastructure.NewUserRepo,
		infrastructure.GetDB,
	)
	return nil
}

func (r *registry) NewAuthUsecase() *usecase.AuthUsecase {
	wire.Build(
		usecase.NewAuthUsecase,
		infrastructure.NewUserRepo,
		infrastructure.GetDB,
	)
	return nil
}
