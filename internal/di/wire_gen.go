// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"auth/internal/api"
	"auth/internal/domain"
	"auth/internal/infrastructure"
	"auth/internal/usecase"
)

// Injectors from wire.go:

func NewServer() *api.Server {
	db := infrastructure.GetDB()
	iUserRepo := infrastructure.NewUserRepo(db)
	iApprovalRepo := infrastructure.NewApprovalRepo(db)
	userService := domain.NewUserService(iUserRepo)
	authUsecase := usecase.NewAuthUsecase(iUserRepo, iApprovalRepo, userService)
	server := api.NewServer(authUsecase)
	return server
}

func NewAuthUsecase() *usecase.AuthUsecase {
	db := infrastructure.GetDB()
	iUserRepo := infrastructure.NewUserRepo(db)
	iApprovalRepo := infrastructure.NewApprovalRepo(db)
	userService := domain.NewUserService(iUserRepo)
	authUsecase := usecase.NewAuthUsecase(iUserRepo, iApprovalRepo, userService)
	return authUsecase
}
