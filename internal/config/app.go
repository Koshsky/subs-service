package config

import (
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
		return nil, err
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
