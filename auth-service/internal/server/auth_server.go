package server

import (
	"context"
	"log"
	"time"

	"github.com/Koshsky/subs-service/auth-service/internal/authpb"
	"github.com/Koshsky/subs-service/auth-service/internal/services"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	AuthService *services.AuthService
}

func NewAuthServer(authService *services.AuthService) *AuthServer {
	return &AuthServer{
		AuthService: authService,
	}
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *authpb.TokenRequest) (*authpb.UserResponse, error) {
	claims, err := s.AuthService.ValidateToken(ctx, req.Token)
	if err != nil {
		return &authpb.UserResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return &authpb.UserResponse{
			Valid: false,
			Error: "Invalid user ID in token",
		}, nil
	}

	email, ok := claims["email"].(string)
	if !ok {
		return &authpb.UserResponse{
			Valid: false,
			Error: "Invalid email in token",
		}, nil
	}

	return &authpb.UserResponse{
		UserId: userIDStr,
		Email:  email,
		Valid:  true,
	}, nil
}

func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	startTime := time.Now()
	log.Printf("[AUTH_SERVER] [%s] Starting Register gRPC handler for email: %s", startTime.Format("15:04:05.000"), req.Email)

	user, err := s.AuthService.Register(ctx, req.Email, req.Password)

	if err != nil {
		totalDuration := time.Since(startTime)
		log.Printf("[AUTH_SERVER] [%s] Register FAILED after %v (service error: %v)", time.Now().Format("15:04:05.000"), totalDuration, err)
		return &authpb.RegisterResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	response := &authpb.RegisterResponse{
		UserId:  user.ID.String(),
		Email:   user.Email,
		Success: true,
		Message: "User created successfully",
	}

	totalDuration := time.Since(startTime)
	log.Printf("[AUTH_SERVER] [%s] Register SUCCESS in %v", time.Now().Format("15:04:05.000"), totalDuration)

	return response, nil
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	token, user, err := s.AuthService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &authpb.LoginResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &authpb.LoginResponse{
		Token:   token,
		UserId:  user.ID.String(),
		Email:   user.Email,
		Success: true,
		Message: "Successful login",
	}, nil
}
