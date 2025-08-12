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

	"github.com/Koshsky/subs-service/notification-service/internal/config"
	"github.com/Koshsky/subs-service/notification-service/internal/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to notification database
	database, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to notification database: %v", err)
	}
	log.Println("Notification database connection established successfully")

	// Ensure proper database connection pool cleanup on shutdown
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to get notification database handle: %v", err)
	}
	defer func() {
		if cerr := sqlDB.Close(); cerr != nil {
			log.Printf("Error closing notification database: %v", cerr)
		}
	}()

	// Initialize RabbitMQ service
	rabbitmqService, err := services.NewRabbitMQService(cfg, database)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ service: %v", err)
	}
	defer rabbitmqService.Close()

	// Start consuming messages
	if err := rabbitmqService.StartConsuming(); err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}

	// Start HTTP server for health checks
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"notification-service","timestamp":"%s"}`,
			time.Now().UTC().Format(time.RFC3339))
	})

	healthSrv := &http.Server{Addr: ":" + cfg.Port, Handler: healthMux}
	go func() {
		log.Printf("Health check server started on port %s", cfg.Port)
		if err := healthSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Health check server failed: %v", err)
		}
	}()

	// Wait for system signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down notification-service...")

	// Context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	// Gracefully stop HTTP health check server
	if err := healthSrv.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown health server: %v", err)
	}

	log.Println("Notification service stopped gracefully")
}
