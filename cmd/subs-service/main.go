package main

import (
	"log"

	"github.com/Koshsky/subs-service/internal/config"
	"github.com/Koshsky/subs-service/internal/middleware"
	"github.com/Koshsky/subs-service/internal/repositories/db"
	"github.com/Koshsky/subs-service/internal/router"
	"github.com/Koshsky/subs-service/internal/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	utils.RegisterCustomValidations()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error while loading app configuration: %v", err)
	}

	db := db.ConnectDatabase(cfg.DB)

	r := gin.New()
	middleware.SetupMiddleware(r, cfg.Middleware)
	router.RegisterRoutes(r, db, cfg.Router)

	log.Println("Starting server on :8080")
	r.Run()
}
