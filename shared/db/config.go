package db

import (
	"fmt"

	"github.com/Koshsky/subs-service/shared/utils"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConnectionString returns the connection string for the database
func (c *DBConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// MigrationDSN returns the connection string for the database for migrations
func (c *DBConfig) MigrationDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

// loadDBConfig loads the database configuration
func loadDBConfig() *DBConfig {
	return &DBConfig{
		Host:     utils.GetEnv("DB_HOST", "postgres"),
		Port:     utils.GetEnv("DB_PORT", "5432"),
		User:     utils.GetEnv("DB_USER", "postgres"),
		Password: utils.GetEnv("DB_PASSWORD", "postgres"),
		DBName:   utils.GetEnv("DB_NAME", "sub-service"),
		SSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
	}
}
