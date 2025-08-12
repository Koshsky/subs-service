package config

import (
	"fmt"

	"github.com/Koshsky/subs-service/auth-service/internal/utils"
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
	Database    DBConfig
	JWTSecret   string
	Port        string
	TLSCertFile string
	TLSKeyFile  string
	EnableTLS   bool
}

func LoadConfig() *Config {
	godotenv.Load()

	db := DBConfig{
		Host:     utils.GetEnv("AUTH_DB_HOST", "auth-db"),
		Port:     utils.GetEnvRequiredWithValidation("AUTH_DB_PORT", utils.ValidatePort),
		User:     utils.GetEnvRequired("AUTH_DB_USER"),
		Password: utils.GetEnvRequired("AUTH_DB_PASSWORD"),
		DBName:   utils.GetEnvRequired("AUTH_DB_NAME"),
		SSLMode:  utils.GetEnv("AUTH_DB_SSLMODE", "disable"),
	}

	return &Config{
		Database:    db,
		JWTSecret:   utils.GetEnvRequiredWithValidation("JWT_SECRET", utils.ValidateMinLength(32)),
		Port:        utils.GetEnvRequiredWithValidation("AUTH_SERVICE_PORT", utils.ValidatePort),
		TLSCertFile: utils.GetEnv("TLS_CERT_FILE", "certs/server-cert.pem"),
		TLSKeyFile:  utils.GetEnv("TLS_KEY_FILE", "certs/server-key.pem"),
		EnableTLS:   utils.GetEnvBool("ENABLE_TLS", false),
	}
}

func (c *Config) ConnectDB() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.Database.ConnectionString()), &gorm.Config{})
}
