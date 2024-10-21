//go:build wireinject

package di

import (
	"auth/internal/api"
	"auth/internal/domain"
	"auth/internal/infrastructure"
	"auth/internal/usecase"

	"github.com/google/wire"
)

func (r *registry) NewServer() *api.Server {
	wire.Build(
		api.NewServer,
		usecase.NewAuthUsecase,
		infrastructure.NewApprovalRepo,
		infrastructure.NewUserRepo,
		domain.NewUserService,
		infrastructure.GetDB,
	)
	return nil
}

func (r *registry) NewAuthUsecase() *usecase.AuthUsecase {
	wire.Build(
		usecase.NewAuthUsecase,
		infrastructure.NewApprovalRepo,
		domain.NewUserService,
		infrastructure.NewUserRepo,
		infrastructure.GetDB,
	)
	return nil
}
