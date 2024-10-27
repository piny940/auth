// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"auth/internal/api"
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
	iClientRepo := gateway.NewClientRepo(db)
	iApprovalRepo := gateway.NewApprovalRepo(db)
	iAuthCodeRepo := gateway.NewAuthCodeRepo(db)
	requestService := oauth.NewRequestService(iClientRepo, iApprovalRepo, iAuthCodeRepo)
	authCodeService := oauth.NewAuthCodeService(iAuthCodeRepo)
	approvalService := oauth.NewApprovalService(iApprovalRepo)
	config := oauth.NewConfig()
	tokenService := oauth.NewTokenService(config, iUserRepo)
	authUsecase := usecase.NewAuthUsecase(userService, requestService, authCodeService, approvalService, tokenService, iUserRepo, iApprovalRepo, iAuthCodeRepo, iClientRepo)
	iClientUsecase := usecase.NewClientUsecase(iClientRepo)
	server := api.NewServer(authUsecase, iClientUsecase)
	return server
}
