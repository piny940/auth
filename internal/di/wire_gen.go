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
	"auth/internal/usecase"
)

// Injectors from wire.go:

func NewServer() *api.Server {
	db := infrastructure.GetDB()
	iUserRepo := infrastructure.NewUserRepo(db)
	iApprovalRepo := infrastructure.NewApprovalRepo(db)
	userService := domain.NewUserService(iUserRepo)
	iClientRepo := infrastructure.NewClientRepo(db)
	authService := oauth.NewAuthService(iClientRepo, iApprovalRepo)
	authUsecase := usecase.NewAuthUsecase(iUserRepo, iApprovalRepo, userService, authService)
	server := api.NewServer(authUsecase)
	return server
}

func NewAuthUsecase() *usecase.AuthUsecase {
	db := infrastructure.GetDB()
	iUserRepo := infrastructure.NewUserRepo(db)
	iApprovalRepo := infrastructure.NewApprovalRepo(db)
	userService := domain.NewUserService(iUserRepo)
	iClientRepo := infrastructure.NewClientRepo(db)
	authService := oauth.NewAuthService(iClientRepo, iApprovalRepo)
	authUsecase := usecase.NewAuthUsecase(iUserRepo, iApprovalRepo, userService, authService)
	return authUsecase
}
