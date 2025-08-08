package main

import (
	"log"

	"github.com/Koshsky/subs-service/core-service/internal/config"
	"github.com/Koshsky/subs-service/core-service/internal/repositories"
	"github.com/Koshsky/subs-service/core-service/internal/router"
	"github.com/Koshsky/subs-service/core-service/internal/services"
)

func main() {
	cfg := config.LoadConfig()

	database, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to core database: %v", err)
	}
	log.Println("Core database connection established successfully")

	authClient, err := services.NewAuthClient(cfg.AuthServiceAddr, cfg.EnableTLS, cfg.TLSCertFile)
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authClient.Close()

	subRepo := repositories.NewSubscriptionRepository(database)
	subService := services.NewSubscriptionService(subRepo)

	r := router.SetupRouter(subService, authClient)

	log.Printf("Core service started on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
