package config

import (
	"fmt"

	"github.com/Koshsky/subs-service/core-service/internal/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (db *DBConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode)
}

type Config struct {
	Database        DBConfig
	Port            string
	AuthServiceAddr string
	TLSCertFile     string
	EnableTLS       bool
}

func LoadConfig() *Config {
	godotenv.Load()

	db := DBConfig{
		Host:     utils.GetEnv("CORE_DB_HOST", "core-db"),
		Port:     utils.GetEnvRequiredWithValidation("CORE_DB_PORT", utils.ValidatePort),
		User:     utils.GetEnvRequired("CORE_DB_USER"),
		Password: utils.GetEnvRequired("CORE_DB_PASSWORD"),
		DBName:   utils.GetEnvRequired("CORE_DB_NAME"),
		SSLMode:  utils.GetEnv("CORE_DB_SSLMODE", "disable"),
	}

	authServicePort := utils.GetEnvRequiredWithValidation("AUTH_SERVICE_PORT", utils.ValidatePort)
	authServiceAddr := "auth-service:" + authServicePort

	return &Config{
		Database:        db,
		Port:            utils.GetEnvRequiredWithValidation("CORE_SERVICE_PORT", utils.ValidatePort),
		AuthServiceAddr: authServiceAddr,
		TLSCertFile:     utils.GetEnv("TLS_CERT_FILE", "certs/server-cert.pem"),
		EnableTLS:       utils.GetEnvBool("ENABLE_TLS", false),
	}
}

func (c *Config) ConnectDB() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.Database.ConnectionString()), &gorm.Config{})
}
