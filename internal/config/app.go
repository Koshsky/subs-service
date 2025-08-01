package config

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DB         *DBConfig
	Router     *RouterConfig
	Middleware *MiddlewareConfig
}

func LoadConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		if _, ok := err.(*fs.PathError); ok {
			// Return both config and a sentinel error that can be checked by caller if needed
			return loadAppConfig(), fmt.Errorf("proceeding without .env file: %w", err)
		}
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}
	return loadAppConfig(), nil
}

func loadAppConfig() *AppConfig {
	return &AppConfig{
		DB:         loadDBConfig(),
		Router:     loadRouterConfig(),
		Middleware: loadMiddlewareConfig(),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return boolValue
	}
	return defaultValue
}
