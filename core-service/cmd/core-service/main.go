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
	// Load application configuration from environment variables
	cfg := config.LoadConfig()

	// Connect to core-service database
	database, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to core database: %v", err)
	}
	log.Println("Core database connection established successfully")

	// Ensure proper database connection pool cleanup on shutdown
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to get core database handle: %v", err)
	}
	defer func() {
		if cerr := sqlDB.Close(); cerr != nil {
			log.Printf("Error closing core database: %v", cerr)
		}
	}()

	// Initialize gRPC client to auth-service (with TLS if enabled)
	authClient, err := services.NewAuthClient(cfg.AuthServiceAddr, cfg.EnableTLS, cfg.TLSCertFile)
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authClient.Close()

	// Set up data access layer and subscription business logic
	subRepo := repositories.NewSubscriptionRepository(database)
	subService := services.NewSubscriptionService(subRepo)

	// Configure HTTP router, middleware and endpoints
	r := router.SetupRouter(subService, authClient, authClient.ValidateToken)

	// Configure HTTP server with timeouts to protect against slow clients
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
		// Hardened and sane defaults to mitigate slowloris and stuck connections
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Start server in a separate goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Printf("Core service started on port %s", cfg.Port)

	// Graceful shutdown: wait for system signals and properly stop the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully stop HTTP server and release resources
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
	log.Println("Server stopped gracefully")
}
