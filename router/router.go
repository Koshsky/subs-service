package router

import (
	"github.com/Koshsky/subs-service/controllers"
	"github.com/Koshsky/subs-service/middleware"
	"github.com/Koshsky/subs-service/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(db *sqlx.DB) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RateLimiterMiddleware())
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
	return r
}
