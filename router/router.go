package router

import (
	"net/http"
	"net/http/pprof"

	"github.com/Koshsky/subs-service/controllers"
	"github.com/Koshsky/subs-service/middleware"
	"github.com/Koshsky/subs-service/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(db *sqlx.DB) *gin.Engine {
	r := gin.New()
	// r.Use(middleware.RateLimiterMiddleware())
	r.Use(middleware.RequestLoggerMiddleware())
	// r.Use(middleware.RequestBodyLoggerMiddleware()) // for debugging
	r.Use(middleware.DatabaseLoggerMiddleware())

	repo := services.NewPostgresRepo(db)
	service := services.NewSubscriptionService(repo)
	controller := controllers.NewSubscriptionController(service)

	r.GET("/subscriptions", controller.List)
	r.POST("/subscriptions", controller.Create)
	r.GET("/subscriptions/:id", controller.Get)
	r.PUT("/subscriptions/:id", controller.Update)
	r.DELETE("/subscriptions/:id", controller.Delete)
	r.GET("/subscriptions/total", controller.SumPrice)

	// Регистрация pprof-обработчиков
	registerPprofHandlers(r)
	return r
}

func registerPprofHandlers(r *gin.Engine) {
	// Группа маршрутов для pprof
	pprofGroup := r.Group("/debug/pprof")
	{
		pprofGroup.GET("/", pprofHandler(pprof.Index))
		pprofGroup.GET("/cmdline", pprofHandler(pprof.Cmdline))
		pprofGroup.GET("/profile", pprofHandler(pprof.Profile))
		pprofGroup.GET("/symbol", pprofHandler(pprof.Symbol))
		pprofGroup.GET("/trace", pprofHandler(pprof.Trace))
		// pprofGroup.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
		// pprofGroup.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
		// pprofGroup.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
		// pprofGroup.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
		// pprofGroup.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
	}
}

// Вспомогательная функция для адаптации pprof-обработчиков к Gin
func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
