package config

import (
	"fmt"

	"github.com/Koshsky/subs-service/auth-service/internal/utils"
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

// Config represents the complete configuration for auth-service
type Config struct {
	Database    DBConfig
	JWTSecret   string
	Port        string
	TLSCertFile string
	TLSKeyFile  string
	EnableTLS   bool
	DatabaseURL string // For backward compatibility
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	godotenv.Load()

	db := DBConfig{
		Host:     utils.GetEnv("AUTH_DB_HOST", "auth-db"),
		Port:     utils.GetEnv("AUTH_DB_PORT", "5432"),
		User:     utils.GetEnv("AUTH_DB_USER", "auth_user"),
		Password: utils.GetEnv("AUTH_DB_PASSWORD", "auth_pass"),
		DBName:   utils.GetEnv("AUTH_DB_NAME", "auth_db"),
		SSLMode:  utils.GetEnv("AUTH_DB_SSLMODE", "disable"),
	}

	return &Config{
		Database:    db,
		DatabaseURL: utils.GetEnv("DATABASE_URL", db.ConnectionString()),
		JWTSecret:   utils.GetEnv("JWT_SECRET", "your_jwt_secret_key"),
		Port:        utils.GetEnv("AUTH_PORT", "50051"),
		TLSCertFile: utils.GetEnv("TLS_CERT_FILE", "certs/server-cert.pem"),
		TLSKeyFile:  utils.GetEnv("TLS_KEY_FILE", "certs/server-key.pem"),
		EnableTLS:   utils.GetEnv("ENABLE_TLS", "true") == "true",
	}
}

// ConnectDB connects to the auth database
func (c *Config) ConnectDB() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.DatabaseURL), &gorm.Config{})
}
