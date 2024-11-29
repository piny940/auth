package server

import (
	"auth/internal/api"
	"auth/internal/di"
	myMiddleware "auth/internal/middleware"
	"context"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var config = &Config{}

type Config struct {
	Port         string   `required:"true"`
	AllowOrigins []string `split_words:"true" required:"true"`
	CSRFEnabled  bool     `split_words:"true" default:"true"`
}

func Init() *echo.Echo {
	err := envconfig.Process("server", config)
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetLevel(log.INFO)
	}
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	if config.CSRFEnabled {
		e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			Skipper:    func(c echo.Context) bool { return strings.HasPrefix(c.Path(), "/oauth") },
			CookiePath: "/",
		}))
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.AllowOrigins,
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderXCSRFToken,
		},
		AllowCredentials: true,
	}))
	e.Use(myMiddleware.Session())
	e.Use(di.NewAuthMiddleware().Auth())
	server := di.NewServer()
	api.RegisterHandlers(e, api.NewStrictHandler(server, nil))
	server.SetLogger(e.Logger)

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
