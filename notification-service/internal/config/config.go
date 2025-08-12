package config

import (
	"fmt"
	"time"

	"github.com/Koshsky/subs-service/notification-service/internal/utils"
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
		Host:     utils.GetEnv("NOTIFY_DB_HOST", "notify-db"),
		Port:     utils.GetEnvRequiredWithValidation("NOTIFY_DB_PORT", utils.ValidatePort),
		User:     utils.GetEnvRequired("NOTIFY_DB_USER"),
		Password: utils.GetEnvRequired("NOTIFY_DB_PASSWORD"),
		DBName:   utils.GetEnvRequired("NOTIFY_DB_NAME"),
		SSLMode:  utils.GetEnv("NOTIFY_DB_SSLMODE", "disable"),
	}

	rabbitmq := RabbitMQConfig{
		URL:      utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		Exchange: utils.GetEnv("RABBITMQ_EXCHANGE", "user_events"),
		Queue:    utils.GetEnv("RABBITMQ_QUEUE", "user_created"),
	}

	shutdownTimeout, _ := time.ParseDuration(utils.GetEnv("NOTIFY_SHUTDOWN_TIMEOUT", "10s"))

	return &Config{
		Database:        db,
		RabbitMQ:        rabbitmq,
		Port:            utils.GetEnv("NOTIFY_SERVICE_PORT", "8082"),
		ShutdownTimeout: shutdownTimeout,
	}
}

// ConnectDB connects to the notification database
func (c *Config) ConnectDB() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.Database.ConnectionString()), &gorm.Config{})
}
