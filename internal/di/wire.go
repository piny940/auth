//go:build wireinject

package di

import (
	"auth/internal/api"
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/gateway"
	"auth/internal/usecase"

	"github.com/google/wire"
)

func (r *registry) NewServer() *api.Server {
	wire.Build(
		api.NewServer,
		usecase.NewAuthUsecase,
		gateway.NewApprovalRepo,
		gateway.NewUserRepo,
		domain.NewUserService,
		oauth.NewAuthService,
		gateway.NewClientRepo,
		infrastructure.GetDB,
	)
	return nil
}

func (r *registry) NewAuthUsecase() *usecase.AuthUsecase {
	wire.Build(
		usecase.NewAuthUsecase,
		gateway.NewApprovalRepo,
		domain.NewUserService,
		gateway.NewUserRepo,
		oauth.NewAuthService,
		gateway.NewClientRepo,
		infrastructure.GetDB,
	)
	return nil
}
