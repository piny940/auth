//go:build !wireinject
// +build !wireinject

package server

import (
	"auth/internal/api"
	"auth/internal/di"
	myMiddleware "auth/internal/middleware"
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var config = &Config{}

type Config struct {
	Port         string   `required:"true"`
	AllowOrigins []string `split_words:"true" required:"true"`
}

func Init() *echo.Echo {
	err := envconfig.Process("server", config)
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.AllowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))
	e.Use(myMiddleware.Session())
	e.Use(myMiddleware.AuthMiddleware())
	api.RegisterHandlers(e, api.NewStrictHandler(di.NewServer(), nil))

	return e
}

func Start(e *echo.Echo) {
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
