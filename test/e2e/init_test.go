//go:build !wireinject
// +build !wireinject

package e2e

import (
	"auth/internal/api"
	"auth/internal/domain/oauth"
	"auth/internal/infrastructure"
	"auth/internal/server"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var baseClient *gorm.DB

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env.test")
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	infrastructure.Init()
	baseClient = infrastructure.GetDB().Client

	code := m.Run()
	os.Exit(code)
}

func newServer(t *testing.T) *httptest.Server {
	t.Cleanup(func() {
		infrastructure.GetDB().Client.Rollback()
	})
	tx := baseClient.Begin()
	infrastructure.InjectDB(&infrastructure.DB{Client: tx})
	api.Init()
	oauth.Init()
	e := server.Init()
	return httptest.NewServer(e)
}
