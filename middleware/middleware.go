package middleware

import (
	"net/http"

	"github.com/Koshsky/subs-service/config"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func SetupMiddleware(r *gin.Engine, cfg *config.MiddlewareConfig) {
	if cfg == nil {
		return
	}

	// Обязательные middleware
	r.Use(gin.Recovery())

	// Опциональные middleware
	if cfg.RateLimiterEnabled {
		r.Use(RateLimiter())
	}
	if cfg.RequestLoggerEnabled {
		r.Use(RequestLogger())
	}
	if cfg.BodyLoggerEnabled {
		r.Use(BodyLogger())
	}
	if cfg.DatabaseLoggerEnabled {
		r.Use(DatabaseLogger())
	}
}

func colorForStatus(code int) string {
	switch {
	case isSuccess(code):
		return color.GreenString("%d", code)
	case isRedirect(code):
		return color.WhiteString("%d", code)
	case isClientError(code):
		return color.YellowString("%d", code)
	default: // server errors
		return color.RedString("%d", code)
	}
}

func isSuccess(code int) bool {
	return code >= http.StatusOK && code < http.StatusMultipleChoices
}

func isRedirect(code int) bool {
	return code >= http.StatusMultipleChoices && code < http.StatusBadRequest
}

func isClientError(code int) bool {
	return code >= http.StatusBadRequest && code < http.StatusInternalServerError
}

func isServerError(code int) bool {
	return code >= http.StatusInternalServerError
}
