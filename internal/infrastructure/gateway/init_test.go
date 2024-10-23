package gateway

import (
	"auth/internal/infrastructure"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var baseClient *gorm.DB

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env.test")
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	infrastructure.Init()
	baseClient = infrastructure.GetDB().Client

	code := m.Run()
	os.Exit(code)
}

func setup(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		infrastructure.GetDB().Client.Rollback()
	})
	tx := baseClient.Begin()
	infrastructure.InjectDB(&infrastructure.DB{Client: tx})
}
