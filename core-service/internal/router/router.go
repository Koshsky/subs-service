package router

import (
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/subs-service/core-service/internal/controllers"
	"github.com/Koshsky/subs-service/core-service/internal/middleware"
)

// SetupRouter sets up the router
func SetupRouter(
	subService controllers.SubscriptionService,
	authClient controllers.AuthClient,
	validateToken middleware.ValidateTokenFunc,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	subController := controllers.NewSubscriptionController(subService)
	authController := controllers.NewAuthController(authClient)

	r := gin.Default()
	r.Use(middleware.RateLimiter())
	r.GET("/health", healthCheck)

	// Auth routes (no auth required)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
	}

	// Protected routes (require authentication)
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(validateToken))
	{
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.POST("", subController.Create)
			subscriptions.GET("", subController.List)
			subscriptions.GET("/:id", subController.Get)
			subscriptions.PUT("/:id", subController.Update)
			subscriptions.DELETE("/:id", subController.Delete)
		}

	}

	// for debugging
	registerPprofHandlers(r)

	return r
}

var pprofPath = "/internal/debug/pprof"

// registerPprofHandlers registers the pprof handlers
func registerPprofHandlers(r *gin.Engine) {
	pprofGroup := r.Group(pprofPath)
	{
		pprofGroup.GET("/", pprofHandler(pprof.Index))
		pprofGroup.GET("/cmdline", pprofHandler(pprof.Cmdline))
		pprofGroup.GET("/profile", pprofHandler(pprof.Profile))
		pprofGroup.GET("/symbol", pprofHandler(pprof.Symbol))
		pprofGroup.GET("/trace", pprofHandler(pprof.Trace))
		pprofGroup.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
		pprofGroup.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
		pprofGroup.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
		pprofGroup.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
		pprofGroup.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
	}
}

// pprofHandler is a helper function to convert a http.HandlerFunc to a gin.HandlerFunc
func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// healthCheck is a health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "core-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
