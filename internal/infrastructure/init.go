package infrastructure

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Client *gorm.DB
}

var db *DB

type Config struct {
	User     string `required:"true"`
	Password string `required:"true"`
	Host     string `required:"true"`
	DBName   string `envconfig:"NAME" required:"true"`
	SSLMode  string `required:"true" envconfig:"DB_SSLMODE"`
	Debug    bool   `default:"false"`
}

func Init() {
	conf := &Config{}
	err := envconfig.Process("DB", conf)
	if err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s",
		conf.User, conf.Password, conf.Host, conf.DBName, conf.SSLMode)
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := client.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	if conf.Debug {
		client = client.Debug()
	}
	db = &DB{Client: client}
}

func GetDB() *DB {
	return db
}

// for test
func InjectDB(new *DB) {
	db = new
}
