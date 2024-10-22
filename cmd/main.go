//go:build !wireinject
// +build !wireinject

package main

import (
	"auth/internal/api"
	"auth/internal/di"
	"auth/internal/infrastructure"
	myMiddleware "auth/internal/middleware"
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Port         string   `required:"true"`
	AllowOrigins []string `split_words:"true" required:"true"`
}

func main() {
	godotenv.Load()
	config := &Config{}
	err := envconfig.Process("server", config)
	if err != nil {
		panic(err)
	}
	infrastructure.Init()
	api.Init()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.AllowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))
	authMid, err := myMiddleware.CreateAuthMiddleware()
	if err != nil {
		panic(err)
	}
	e.Use(authMid)
	api.RegisterHandlers(e, api.NewStrictHandler(di.NewServer(), nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(":" + config.Port); err != nil {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
