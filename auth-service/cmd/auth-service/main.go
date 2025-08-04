package main

import (
	"log"
	"net"

	"github.com/Koshsky/subs-service/auth-service/internal/authpb"
	"github.com/Koshsky/subs-service/auth-service/internal/config"
	"github.com/Koshsky/subs-service/auth-service/internal/repositories"
	"github.com/Koshsky/subs-service/auth-service/internal/server"
	"github.com/Koshsky/subs-service/auth-service/internal/services"
	"github.com/Koshsky/subs-service/shared/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.LoadConfig()

	db := db.ConnectDatabase(cfg.DatabaseURL)

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, []byte(cfg.JWTSecret))
	authServer := server.NewAuthServer(authService)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var grpcServer *grpc.Server

	if cfg.EnableTLS {
		// Create TLS credentials
		creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertFile, cfg.TLSKeyFile)
		if err != nil {
			log.Printf("Failed to load TLS credentials, starting without TLS: %v", err)
			grpcServer = grpc.NewServer()
			log.Printf("Auth service started without TLS on port %s (WARNING: Insecure)", cfg.Port)
		} else {
			grpcServer = grpc.NewServer(grpc.Creds(creds))
			log.Printf("Auth service started with TLS on port %s", cfg.Port)
		}
	} else {
		grpcServer = grpc.NewServer()
		log.Printf("Auth service started without TLS on port %s (WARNING: Insecure)", cfg.Port)
	}

	authpb.RegisterAuthServiceServer(grpcServer, authServer)

	// Включаем reflection для отладки
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
