// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"auth/internal/api"
	"auth/internal/api/middleware"
	"auth/internal/domain"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/gateway"
	"auth/internal/usecase"
)

// Injectors from wire.go:

func NewServer() *api.Server {
	db := infrastructure.GetDB()
	iUserRepo := gateway.NewUserRepo(db)
	userService := domain.NewUserService(iUserRepo)
	userUsecase := usecase.NewAuthUsecase(userService, iUserRepo)
	iAuthCodeRepo := gateway.NewAuthCodeRepo(db)
	authCodeService := oauth.NewAuthCodeService(iAuthCodeRepo)
	config := oauth.NewConfig()
	jwKsService := oauth.NewJWKsService(config)
	iApprovalRepo := gateway.NewApprovalRepo(db)
	approvalService := oauth.NewApprovalService(iApprovalRepo)
	tokenService := oauth.NewTokenService(config, iUserRepo)
	iClientRepo := gateway.NewClientRepo(db)
	oAuthUsecase := usecase.NewOAuthUsecase(authCodeService, jwKsService, approvalService, iApprovalRepo, tokenService, iClientRepo)
	iClientUsecase := usecase.NewClientUsecase(iClientRepo)
	approvalUsecase := usecase.NewApprovalUsecase(iApprovalRepo, iClientRepo)
	middlewareConfig := middleware.NewConfig()
	echoContextReg := middleware.NewEchoContextReg()
	auth := middleware.NewAuth(middlewareConfig, echoContextReg)
	server := api.NewServer(userUsecase, oAuthUsecase, iClientUsecase, approvalUsecase, auth)
	return server
}

func NewAuthMiddleware() *middleware.AuthMiddleware {
	config := middleware.NewConfig()
	db := infrastructure.GetDB()
	iClientRepo := gateway.NewClientRepo(db)
	iUserRepo := gateway.NewUserRepo(db)
	authMiddleware := middleware.NewAuthMiddleware(config, iClientRepo, iUserRepo)
	return authMiddleware
}

func NewEchoContextMiddleware() *middleware.EchoContextMiddleware {
	echoContextMiddleware := middleware.NewEchoContextMiddleware()
	return echoContextMiddleware
}

func NewTokenService() *oauth.TokenService {
	config := oauth.NewConfig()
	db := infrastructure.GetDB()
	iUserRepo := gateway.NewUserRepo(db)
	tokenService := oauth.NewTokenService(config, iUserRepo)
	return tokenService
}
