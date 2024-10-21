//go:build wireinject

package di

import (
	"auth/internal/api"
	"auth/internal/domain"
	"auth/internal/domain/oauth"
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
		oauth.NewAuthService,
		infrastructure.NewClientRepo,
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
		oauth.NewAuthService,
		infrastructure.NewClientRepo,
		infrastructure.GetDB,
	)
	return nil
}
