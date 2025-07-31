package router

import (
	"net/http"
	"net/http/pprof"

	"github.com/Koshsky/subs-service/internal/config"
	"github.com/Koshsky/subs-service/internal/controllers"
	"github.com/Koshsky/subs-service/internal/repositories/sub_repository"
	"github.com/Koshsky/subs-service/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var pprofPath = "/internal/debug/pprof"

func RegisterRoutes(r *gin.Engine, conn *gorm.DB, cfg *config.RouterConfig) {
	r.GET("/health", healthCheck)

	registerSubHandlers(r, conn)

	if cfg.EnableProfiling {
		registerPprofHandlers(r)
	}
}

func registerSubHandlers(r *gin.Engine, conn *gorm.DB) {
	repo := sub_repository.New(conn)
	service := services.NewSubService(repo)
	subController := controllers.NewSubscriptionController(service)

	r.GET("/subscriptions", subController.List)
	r.POST("/subscriptions", subController.Create)
	r.GET("/subscriptions/:id", subController.Get)
	r.PUT("/subscriptions/:id", subController.Update)
	r.DELETE("/subscriptions/:id", subController.Delete)
	r.GET("/subscriptions/total", subController.SumPrice)

}

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

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
