package server

import (
	"auth/internal/api"
	"auth/internal/api/middleware"
	"auth/internal/di"
	"context"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	echoMid "github.com/labstack/echo/v4/middleware"
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
	e.Use(echoMid.RequestID())
	e.Use(echoMid.Logger())
	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetLevel(log.INFO)
	}
	e.Use(echoMid.Recover())
	e.Use(echoMid.Secure())
	if config.CSRFEnabled {
		e.Use(echoMid.CSRFWithConfig(echoMid.CSRFConfig{
			Skipper:    func(c echo.Context) bool { return strings.HasPrefix(c.Path(), "/oauth") },
			CookiePath: "/",
		}))
	}
	e.Use(echoMid.CORSWithConfig(echoMid.CORSConfig{
		AllowOrigins: config.AllowOrigins,
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderXCSRFToken,
		},
		AllowCredentials: true,
	}))
	e.Use(di.NewAuthMiddleware().Auth())
	e.Use(middleware.NewEchoContextMiddleware().Context())
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
