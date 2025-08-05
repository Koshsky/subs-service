package config

import (
	"fmt"

	"github.com/Koshsky/subs-service/shared/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConfig represents database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConnectionString returns the connection string for the database
func (db *DBConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode)
}

// Config represents the complete configuration for core-service
type Config struct {
	Database        DBConfig
	Port            string
	AuthServiceAddr string
	TLSCertFile     string
	EnableTLS       bool
	DatabaseURL     string // For backward compatibility
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	godotenv.Load()

	db := DBConfig{
		Host:     utils.GetEnv("CORE_DB_HOST", "core-db"),
		Port:     utils.GetEnv("CORE_DB_PORT", "5432"),
		User:     utils.GetEnv("CORE_DB_USER", "core_user"),
		Password: utils.GetEnv("CORE_DB_PASSWORD", "core_pass"),
		DBName:   utils.GetEnv("CORE_DB_NAME", "core_db"),
		SSLMode:  utils.GetEnv("CORE_DB_SSLMODE", "disable"),
	}

	return &Config{
		Database:        db,
		DatabaseURL:     utils.GetEnv("DATABASE_URL", db.ConnectionString()),
		Port:            utils.GetEnv("CORE_PORT", "8080"),
		AuthServiceAddr: utils.GetEnv("AUTH_SERVICE_ADDR", "localhost:50051"),
		TLSCertFile:     utils.GetEnv("TLS_CERT_FILE", "certs/server-cert.pem"),
		EnableTLS:       utils.GetEnv("ENABLE_TLS", "true") == "true",
	}
}

// ConnectDB connects to the core database
func (c *Config) ConnectDB() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.DatabaseURL), &gorm.Config{})
}
