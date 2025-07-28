package config

type RouterConfig struct {
	Middleware      MiddlewareConfig
	EnableProfiling bool
	// Другие настройки роутера при необходимости
	// Timeout int
	// MaxBodySize int64
}

type MiddlewareConfig struct {
	RateLimiterEnabled    bool
	RequestLoggerEnabled  bool
	BodyLoggerEnabled     bool
	DatabaseLoggerEnabled bool
}

func loadRouterConfig() *RouterConfig {
	return &RouterConfig{
		Middleware: MiddlewareConfig{
			RateLimiterEnabled:    getEnvBool("RATE_LIMITER_ENABLED", false),
			RequestLoggerEnabled:  getEnvBool("REQUEST_LOGGER_ENABLED", true),
			BodyLoggerEnabled:     getEnvBool("BODY_LOGGER_ENABLED", false),
			DatabaseLoggerEnabled: getEnvBool("DB_LOGGER_ENABLED", false),
		},
		EnableProfiling: getEnvBool("ENABLE_PROFILING", false),
	}
}
