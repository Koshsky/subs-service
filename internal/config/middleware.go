package config

type MiddlewareConfig struct {
	RateLimiterEnabled    bool
	RequestLoggerEnabled  bool
	BodyLoggerEnabled     bool
	DatabaseLoggerEnabled bool
}

func loadMiddlewareConfig() *MiddlewareConfig {
	return &MiddlewareConfig{
		RateLimiterEnabled:    getEnvBool("RATE_LIMITER_ENABLED", false),
		RequestLoggerEnabled:  getEnvBool("REQUEST_LOGGER_ENABLED", true),
		BodyLoggerEnabled:     getEnvBool("BODY_LOGGER_ENABLED", false),
		DatabaseLoggerEnabled: getEnvBool("DB_LOGGER_ENABLED", false),
	}
}
