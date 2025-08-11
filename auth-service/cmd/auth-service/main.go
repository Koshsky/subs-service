package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Koshsky/subs-service/auth-service/internal/authpb"
	"github.com/Koshsky/subs-service/auth-service/internal/config"
	"github.com/Koshsky/subs-service/auth-service/internal/repositories"
	"github.com/Koshsky/subs-service/auth-service/internal/server"
	"github.com/Koshsky/subs-service/auth-service/internal/services"
)

func main() {
	// Load authentication service configuration
	cfg := config.LoadConfig()

	// Connect to auth-service database
	database, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to auth database: %v", err)
	}
	log.Println("Auth database connection established successfully")

	// Ensure proper database connection pool cleanup on shutdown
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to get auth database handle: %v", err)
	}
	defer func() {
		if cerr := sqlDB.Close(); cerr != nil {
			log.Printf("Error closing auth database: %v", cerr)
		}
	}()

	// Initialize repository and authentication business logic
	userRepo := repositories.NewUserRepository(database)
	authService := services.NewAuthService(userRepo, []byte(cfg.JWTSecret))
	authServer := server.NewAuthServer(authService)

	// Start HTTP server for health checks with graceful shutdown capability
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"auth-service","timestamp":"%s"}`,
			time.Now().UTC().Format(time.RFC3339))
	})
	healthPort := "8081" // Fixed health check port
	healthSrv := &http.Server{Addr: ":" + healthPort, Handler: healthMux}
	go func() {
		log.Printf("Health check server started on port %s", healthPort)
		if err := healthSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Health check server failed: %v", err)
		}
	}()

	// Create and configure gRPC server
	grpcServer, err := server.NewGRPCServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create gRPC server: %v", err)
	}

	// Register gRPC handlers for authentication service
	authpb.RegisterAuthServiceServer(grpcServer.GetServer(), authServer)

	// Start gRPC server in a separate goroutine to allow graceful shutdown
	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	// Wait for system signals for graceful service shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down auth-service...")

	// Context with timeout, shared for stopping gRPC and health server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1) Gracefully stop gRPC server
	grpcServer.Stop()

	// 2) Gracefully stop HTTP health check server
	if err := healthSrv.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown health server: %v", err)
	}
}
