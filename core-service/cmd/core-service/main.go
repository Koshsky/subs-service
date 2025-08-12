package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to get core database handle: %v", err)
	}
	defer func() {
		if cerr := sqlDB.Close(); cerr != nil {
			log.Printf("Error closing core database: %v", cerr)
		}
	}()

	authClient, err := services.NewAuthClient(cfg.AuthServiceAddr, cfg.EnableTLS, cfg.TLSCertFile)
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authClient.Close()

	subRepo := repositories.NewSubscriptionRepository(database)
	subService := services.NewSubscriptionService(subRepo)

	r := router.SetupRouter(subService, authClient, authClient.ValidateToken)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
	log.Println("Server stopped gracefully")
}
