//go:build !wireinject
// +build !wireinject

package main

import (
	"auth/internal/api"
	"auth/internal/infrastructure"
	"auth/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	infrastructure.Init()
	api.Init()
	e := server.Init()

	server.Start(e)
}
