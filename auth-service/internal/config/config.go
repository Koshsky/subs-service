package config

import (
	"github.com/Koshsky/subs-service/shared/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		DatabaseURL: utils.GetEnv("DATABASE_URL", "postgres://user:password@localhost/subs_db?sslmode=disable"),
		JWTSecret:   utils.GetEnv("JWT_SECRET", "your_jwt_secret_key"),
		Port:        utils.GetEnv("AUTH_PORT", "50051"),
	}
}

func (c *Config) ConnectDB() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.DatabaseURL), &gorm.Config{})
}
