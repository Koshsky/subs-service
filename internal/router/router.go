package router

import (
	"net/http"
	"net/http/pprof"

	"github.com/Koshsky/subs-service/internal/config"
	"github.com/Koshsky/subs-service/internal/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var pprofPath = "/internal/debug/pprof"

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.RouterConfig) {

	r.GET("/subscriptions", controllers.List(db))
	r.POST("/subscriptions", controllers.Create(db))
	r.GET("/subscriptions/:id", controllers.Get(db))
	r.PUT("/subscriptions/:id", controllers.Update(db))
	r.DELETE("/subscriptions/:id", controllers.Delete(db))
	r.GET("/subscriptions/total", controllers.SumPrice(db))
	r.GET("/health", healthCheck)

	if cfg.EnableProfiling {
		registerPprofHandlers(r)
	}
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
