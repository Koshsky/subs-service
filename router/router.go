package router

import (
	"github.com/Koshsky/subs-service/controllers"
	"github.com/Koshsky/subs-service/middleware"
	"github.com/Koshsky/subs-service/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	r.Use(middleware.RateLimiterMiddleware())
	r.Use(middleware.RequestLoggerMiddleware())
	// r.Use(middleware.RequestBodyLoggerMiddleware()) // for debugging
	r.Use(middleware.DatabaseLoggerMiddleware())

	service := services.NewSubscriptionService(db)
	controller := controllers.NewSubscriptionController(service)

	r.GET("/subscriptions", controller.List)
	r.POST("/subscriptions", controller.Create)
	r.GET("/subscriptions/:id", controller.Get)
	r.PUT("/subscriptions/:id", controller.Update)
	r.DELETE("/subscriptions/:id", controller.Delete)
	r.GET("/subscriptions/total", controller.SumPrice)
}
