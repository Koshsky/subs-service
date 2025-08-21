package main

import (
	"log"
	"net"

	"github.com/Koshsky/subs-service/auth-service/internal/authpb"
	"github.com/Koshsky/subs-service/auth-service/internal/config"
	"github.com/Koshsky/subs-service/auth-service/internal/messaging"
	"github.com/Koshsky/subs-service/auth-service/internal/repositories"
	"github.com/Koshsky/subs-service/auth-service/internal/server"
	"github.com/Koshsky/subs-service/auth-service/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize RabbitMQ service
	rabbitmqService, err := messaging.NewRabbitMQAdapter(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize RabbitMQ service: %v", err)
		log.Printf("Auth service will continue without event publishing")
		rabbitmqService = nil
	} else {
		defer rabbitmqService.Close()
	}

	gormAdapter, err := repositories.NewGormAdapter(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create database adapter: %v", err)
	}
	userRepo := repositories.NewUserRepository(gormAdapter)
	authService := services.NewAuthService(userRepo, rabbitmqService, cfg)
	authServer := server.NewAuthServer(authService)

	var grpcServer *grpc.Server
	if cfg.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertFile, cfg.TLSKeyFile)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds))
	} else {
		grpcServer = grpc.NewServer()
	}

	authpb.RegisterAuthServiceServer(grpcServer, authServer)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("Auth service starting on port %s", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("gRPC server stopped: %v", err)
	}
}
