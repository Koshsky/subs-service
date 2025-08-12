package services

import (
	"context"
	"log"

	"github.com/Koshsky/subs-service/core-service/internal/corepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client corepb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(authServiceAddr string, enableTLS bool, tlsCertFile string) (*AuthClient, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)

	if enableTLS {
		// Server name must match the certificate (SAN/CN)
		creds, cerr := credentials.NewClientTLSFromFile(tlsCertFile, "")
		if cerr != nil {
			return nil, cerr
		}
		conn, err = grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(creds))
	} else {
		conn, err = grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if err != nil {
		return nil, err
	}

	client := corepb.NewAuthServiceClient(conn)

	return &AuthClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the connection to the auth service
func (ac *AuthClient) Close() {
	if ac.conn != nil {
		ac.conn.Close()
	}
}

// ValidateToken validates the token
func (ac *AuthClient) ValidateToken(ctx context.Context, token string) (*corepb.UserResponse, error) {
	req := &corepb.TokenRequest{Token: token}
	resp, err := ac.client.ValidateToken(ctx, req)
	if err != nil {
		log.Printf("Failed to validate token: %v", err)
		return nil, err
	}
	return resp, nil
}

// Register registers a new user
func (ac *AuthClient) Register(ctx context.Context, email, password string) (*corepb.RegisterResponse, error) {
	req := &corepb.RegisterRequest{Email: email, Password: password}
	resp, err := ac.client.Register(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Login logs in a user
func (ac *AuthClient) Login(ctx context.Context, email, password string) (*corepb.LoginResponse, error) {
	req := &corepb.LoginRequest{Email: email, Password: password}
	resp, err := ac.client.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
