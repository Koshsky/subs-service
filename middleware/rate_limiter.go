package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

const (
	RequestsPerSecond = 1
	BurstLimit        = 2 // Допускаем кратковременные всплески
)

func RateLimiter() gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(RequestsPerSecond, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Minute,
	})

	lmt.SetBurst(BurstLimit)
	lmt.SetMessage("Rate limit exceeded")
	lmt.SetMessageContentType("application/json")

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		path := c.Request.URL.Path
		method := c.Request.Method

		httpErr := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpErr != nil {
			retryAfter := int(time.Second) * BurstLimit

			log.Printf("%s %s %s %s %s %s",
				color.CyanString("[RATE]"),
				color.RedString("BLOCKED"),
				color.MagentaString(clientIP),
				color.YellowString(method),
				path,
				color.RedString("(wait %v)", time.Duration(retryAfter)),
			)

			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       httpErr.Message,
				"retry_after": retryAfter,
				"rate_limit":  RequestsPerSecond,
				"burst":       BurstLimit,
			})
			return
		}

		log.Printf("%s %s %s %s %s",
			color.CyanString("[RATE]"),
			color.GreenString("ALLOWED"),
			color.MagentaString(clientIP),
			color.YellowString(method),
			path,
		)

		c.Next()
	}
}
