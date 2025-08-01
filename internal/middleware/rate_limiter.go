package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	RequestsPerSecond = 10
	BurstLimit        = 20
	CleanupInterval   = 5 * time.Minute  // How often to clean up old entries
	IPExpiration      = 10 * time.Minute // When to remove inactive IPs
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	limiters = make(map[string]*ipLimiter)
	mu       sync.RWMutex
	once     sync.Once
)

func RateLimiter() gin.HandlerFunc {
	// Start cleanup goroutine just once
	once.Do(func() {
		go cleanupExpiredIPs()
	})

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Get or create the limiter for this IP
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too Many Requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	// Update last seen time if IP exists
	if limiter, exists := limiters[ip]; exists {
		limiter.lastSeen = time.Now()
		return limiter.limiter
	}

	// Create new limiter for new IP
	limiters[ip] = &ipLimiter{
		limiter:  rate.NewLimiter(RequestsPerSecond, BurstLimit),
		lastSeen: time.Now(),
	}

	return limiters[ip].limiter
}

func cleanupExpiredIPs() {
	for {
		time.Sleep(CleanupInterval)

		mu.Lock()
		for ip, limiter := range limiters {
			if time.Since(limiter.lastSeen) > IPExpiration {
				delete(limiters, ip)
			}
		}
		mu.Unlock()
	}
}
