//go:build wireinject

package di

import (
	"auth/internal/api"
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/gateway"
	"auth/internal/middleware"
	"auth/internal/usecase"

	"github.com/google/wire"
)

func NewServer() *api.Server {
	wire.Build(
		api.NewServer,
		usecase.NewAuthUsecase,
		usecase.NewOAuthUsecase,
		gateway.NewApprovalRepo,
		gateway.NewUserRepo,
		domain.NewUserService,
		oauth.NewApprovalService,
		oauth.NewAuthCodeService,
		oauth.NewRequestService,
		oauth.NewTokenService,
		oauth.NewConfig,
		gateway.NewClientRepo,
		gateway.NewAuthCodeRepo,
		usecase.NewClientUsecase,
		infrastructure.GetDB,
	)
	return nil
}

func NewAuthMiddleware() *middleware.AuthMiddleware {
	wire.Build(
		middleware.NewAuthMiddleware,
		gateway.NewClientRepo,
		infrastructure.GetDB,
	)
	return nil
}
