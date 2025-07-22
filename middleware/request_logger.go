package middleware

import (
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Printf("%s %s %s %s %s",
			color.CyanString("[GIN]"),
			color.YellowString("%-7s", c.Request.Method),
			c.Request.URL.Path,
			colorForStatus(c.Writer.Status()),
			duration,
		)
	}
}
