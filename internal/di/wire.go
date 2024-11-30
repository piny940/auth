//go:build wireinject

package di

import (
	"auth/internal/api"
	"auth/internal/api/middleware"
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/gateway"
	"auth/internal/usecase"

	"github.com/google/wire"
)

func NewServer() *api.Server {
	wire.Build(
		api.NewServer,
		middleware.NewAuth,
		middleware.NewConfig,
		middleware.NewEchoContextReg,
		usecase.NewAuthUsecase,
		usecase.NewOAuthUsecase,
		gateway.NewApprovalRepo,
		gateway.NewUserRepo,
		domain.NewUserService,
		oauth.NewApprovalService,
		oauth.NewAuthCodeService,
		oauth.NewTokenService,
		oauth.NewJWKsService,
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
		middleware.NewConfig,
		middleware.NewAuthMiddleware,
		gateway.NewUserRepo,
		gateway.NewClientRepo,
		infrastructure.GetDB,
	)
	return nil
}

func NewEchoContextMiddleware() *middleware.EchoContextMiddleware {
	wire.Build(
		middleware.NewEchoContextMiddleware,
	)
	return nil
}

func NewTokenService() *oauth.TokenService {
	wire.Build(
		oauth.NewTokenService,
		oauth.NewConfig,
		gateway.NewUserRepo,
		infrastructure.GetDB,
	)
	return nil
}
