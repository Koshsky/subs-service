package main

import (
	"log"
	"net"

	"github.com/Koshsky/subs-service/auth-service/internal/authpb"
	"github.com/Koshsky/subs-service/auth-service/internal/config"
	"github.com/Koshsky/subs-service/auth-service/internal/repositories"
	"github.com/Koshsky/subs-service/auth-service/internal/server"
	"github.com/Koshsky/subs-service/auth-service/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cfg := config.LoadConfig()

	database, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to auth database: %v", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to get auth database handle: %v", err)
	}
	defer func() {
		if cerr := sqlDB.Close(); cerr != nil {
			log.Printf("Error closing auth database: %v", cerr)
		}
	}()

	userRepo := repositories.NewUserRepository(database)
	authService := services.NewAuthService(userRepo, []byte(cfg.JWTSecret))
	authServer := server.NewAuthServer(authService)

	var grpcServer *grpc.Server
	if cfg.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertFile, cfg.TLSKeyFile)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds))
		log.Printf("Auth service configured with TLS")
	} else {
		grpcServer = grpc.NewServer()
		log.Printf("Auth service configured without TLS (WARNING: Insecure)")
	}

	authpb.RegisterAuthServiceServer(grpcServer, authServer)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("Auth service starting on port %s", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("gRPC server stopped: %v", err)
	}
}
