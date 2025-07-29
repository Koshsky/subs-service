package router

import (
	"net/http"
	"net/http/pprof"

	"github.com/Koshsky/subs-service/internal/config"
	"github.com/Koshsky/subs-service/internal/controllers"
	"github.com/Koshsky/subs-service/internal/middleware"
	"github.com/Koshsky/subs-service/internal/repositories"
	"github.com/Koshsky/subs-service/internal/services"
	"github.com/gin-gonic/gin"
)

var pprofPath = "/internal/debug/pprof"

func SetupRouter(repo repositories.SubscriptionRepository, routerCfg *config.RouterConfig) *gin.Engine {
	r := gin.New()

	middleware.SetupMiddleware(r, &routerCfg.Middleware)

	controller := initController(repo)

	r.GET("/subscriptions", controller.List)
	r.POST("/subscriptions", controller.Create)
	r.GET("/subscriptions/:id", controller.Get)
	r.PUT("/subscriptions/:id", controller.Update)
	r.DELETE("/subscriptions/:id", controller.Delete)
	r.GET("/subscriptions/total", controller.SumPrice)
	r.GET("/health", healthCheck)

	if routerCfg.EnableProfiling {
		registerPprofHandlers(r)
	}
	return r
}

func initController(repo repositories.SubscriptionRepository) *controllers.SubscriptionController {
	service := services.New(repo)
	return controllers.New(service)
}

func registerPprofHandlers(r *gin.Engine) {
	pprofGroup := r.Group(pprofPath)
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

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
