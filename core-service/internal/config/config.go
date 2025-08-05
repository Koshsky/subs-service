package config

import (
	"github.com/Koshsky/subs-service/shared/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DatabaseURL     string
	Port            string
	AuthServiceAddr string
	TLSCertFile     string
	EnableTLS       bool
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		DatabaseURL:     utils.GetEnv("DATABASE_URL", "postgres://user:password@localhost/subs_db?sslmode=disable"),
		Port:            utils.GetEnv("CORE_PORT", "8080"),
		AuthServiceAddr: utils.GetEnv("AUTH_SERVICE_ADDR", "localhost:50051"),
		TLSCertFile:     utils.GetEnv("TLS_CERT_FILE", "certs/server-cert.pem"),
		EnableTLS:       utils.GetEnv("ENABLE_TLS", "true") == "true",
	}
}

func (c *Config) ConnectDB() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.DatabaseURL), &gorm.Config{})
}
