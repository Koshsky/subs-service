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

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, authServer)

	// Включаем reflection для отладки
	reflection.Register(grpcServer)

	log.Printf("Auth service started on port %s", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
