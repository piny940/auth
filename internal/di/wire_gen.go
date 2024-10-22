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
	iApprovalRepo := gateway.NewApprovalRepo(db)
	userService := domain.NewUserService(iUserRepo)
	iClientRepo := gateway.NewClientRepo(db)
	authService := oauth.NewAuthService(iClientRepo, iApprovalRepo)
	authUsecase := usecase.NewAuthUsecase(iUserRepo, iApprovalRepo, userService, authService)
	iClientUsecase := usecase.NewClientUsecase(iClientRepo)
	server := api.NewServer(authUsecase, iClientUsecase)
	return server
}
