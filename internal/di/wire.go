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
		gateway.NewAuthCodeRepo,
		usecase.NewClientUsecase,
		infrastructure.GetDB,
	)
	return nil
}
