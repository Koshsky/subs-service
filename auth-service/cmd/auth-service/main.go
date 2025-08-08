package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Koshsky/subs-service/auth-service/internal/authpb"
	"github.com/Koshsky/subs-service/auth-service/internal/config"
	"github.com/Koshsky/subs-service/auth-service/internal/repositories"
	"github.com/Koshsky/subs-service/auth-service/internal/server"
	"github.com/Koshsky/subs-service/auth-service/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	cfg := config.LoadConfig()

	database, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to auth database: %v", err)
	}
	log.Println("Auth database connection established successfully")

	userRepo := repositories.NewUserRepository(database)
	authService := services.NewAuthService(userRepo, []byte(cfg.JWTSecret))
	authServer := server.NewAuthServer(authService)

	// Start HTTP health check server
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"status":"ok","service":"auth-service","timestamp":"%s"}`,
				time.Now().UTC().Format(time.RFC3339))
		})

		healthPort := "8081" // Fixed health check port
		log.Printf("Health check server started on port %s", healthPort)
		if err := http.ListenAndServe(":"+healthPort, nil); err != nil {
			log.Printf("Health check server failed: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var grpcServer *grpc.Server

	if cfg.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertFile, cfg.TLSKeyFile)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials for auth-service: %v", err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds))
		log.Printf("Auth service started with TLS on port %s", cfg.Port)
	} else {
		grpcServer = grpc.NewServer()
		log.Printf("Auth service started without TLS on port %s (WARNING: Insecure)", cfg.Port)
	}

	authpb.RegisterAuthServiceServer(grpcServer, authServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
