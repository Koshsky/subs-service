package services

import (
	"context"
	"log"
	"time"

	"github.com/Koshsky/subs-service/core-service/internal/corepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
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

	dialOptions := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             1 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024),
			grpc.MaxCallSendMsgSize(1024*1024),
		),
	}

	if enableTLS {
		creds, cerr := credentials.NewClientTLSFromFile(tlsCertFile, "")
		if cerr != nil {
			return nil, cerr
		}
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))
	} else {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err = grpc.NewClient(authServiceAddr, dialOptions...)
	if err != nil {
		return nil, err
	}

	client := corepb.NewAuthServiceClient(conn)

	ac := &AuthClient{
		client: client,
		conn:   conn,
	}

	go ac.warmupConnection()

	return ac, nil
}

func (ac *AuthClient) Close() {
	if ac.conn != nil {
		ac.conn.Close()
	}
}

func (ac *AuthClient) warmupConnection() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	req := &corepb.TokenRequest{Token: ""}
	_, err := ac.client.ValidateToken(ctx, req)
	if err != nil {
		log.Printf("Warmup connection completed (expected error: %v)", err)
	} else {
		log.Printf("Warmup connection completed successfully")
	}
}

func (ac *AuthClient) ValidateToken(ctx context.Context, token string) (*corepb.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	req := &corepb.TokenRequest{Token: token}
	resp, err := ac.client.ValidateToken(ctx, req)
	if err != nil {
		log.Printf("Failed to validate token: %v", err)
		return nil, err
	}
	return resp, nil
}

func (ac *AuthClient) Register(ctx context.Context, email, password string) (*corepb.RegisterResponse, error) {
	req := &corepb.RegisterRequest{Email: email, Password: password}
	resp, err := ac.client.Register(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ac *AuthClient) Login(ctx context.Context, email, password string) (*corepb.LoginResponse, error) {
	req := &corepb.LoginRequest{Email: email, Password: password}
	resp, err := ac.client.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
