package server

import (
	"fmt"
	"log"
	"net"

	"github.com/Koshsky/subs-service/auth-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCServer struct {
	server *grpc.Server
	config *config.Config
}

// NewGRPCServer creates a new gRPC server with optional TLS support
func NewGRPCServer(cfg *config.Config) (*GRPCServer, error) {
	var grpcServer *grpc.Server

	// Enable TLS for gRPC if configured
	if cfg.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertFile, cfg.TLSKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create TLS credentials for auth-service: %v", err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds))
		log.Printf("Auth service configured with TLS")
	} else {
		grpcServer = grpc.NewServer()
		log.Printf("Auth service configured without TLS (WARNING: Insecure)")
	}

	return &GRPCServer{
		server: grpcServer,
		config: cfg,
	}, nil
}

// GetServer returns the underlying gRPC server instance
func (g *GRPCServer) GetServer() *grpc.Server {
	return g.server
}

// Start starts the gRPC server
func (g *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", ":"+g.config.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", g.config.Port, err)
	}

	if g.config.EnableTLS {
		log.Printf("Auth service starting with TLS on port %s", g.config.Port)
	} else {
		log.Printf("Auth service starting without TLS on port %s (WARNING: Insecure)", g.config.Port)
	}

	return g.server.Serve(lis)
}

// Stop gracefully stops the gRPC server
func (g *GRPCServer) Stop() {
	if g.server != nil {
		g.server.GracefulStop()
		log.Println("Auth service stopped gracefully")
	}
}
