package config

import (
	"fmt"
	"os"
	"time"

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

// RabbitMQConfig represents RabbitMQ configuration
type RabbitMQConfig struct {
	URL      string
	Exchange string
	Queue    string
}

// Config represents the complete configuration for notification-service
type Config struct {
	Database        DBConfig
	RabbitMQ        RabbitMQConfig
	Port            string
	ShutdownTimeout time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	godotenv.Load()

	db := DBConfig{
		Host:     getEnv("NOTIFY_DB_HOST", ""),
		Port:     getEnv("NOTIFY_DB_PORT", ""),
		User:     getEnv("NOTIFY_DB_USER", ""),
		Password: getEnv("NOTIFY_DB_PASSWORD", ""),
		DBName:   getEnv("NOTIFY_DB_NAME", ""),
		SSLMode:  getEnv("NOTIFY_DB_SSLMODE", ""),
	}

	rabbitmq := RabbitMQConfig{
		URL:      getEnv("RABBITMQ_URL", ""),
		Exchange: getEnv("RABBITMQ_EXCHANGE", ""),
		Queue:    getEnv("RABBITMQ_QUEUE", ""),
	}

	shutdownTimeout, _ := time.ParseDuration(getEnv("NOTIFY_SHUTDOWN_TIMEOUT", "10s"))

	return &Config{
		Database:        db,
		RabbitMQ:        rabbitmq,
		Port:            getEnv("NOTIFY_SERVICE_PORT", ""),
		ShutdownTimeout: shutdownTimeout,
	}
}

// ConnectDB connects to the notification database
func (c *Config) ConnectDB() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.Database.ConnectionString()), &gorm.Config{})
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
