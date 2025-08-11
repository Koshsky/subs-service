#!/bin/bash

# Script to generate protobuf files
set -e

echo "🔧 Generating protobuf files..."

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc not found. Install Protocol Buffers compiler:"
    echo "   Ubuntu/Debian: sudo apt install protobuf-compiler"
    echo "   macOS: brew install protobuf"
    exit 1
fi

# Check if protoc-gen-go and protoc-gen-go-grpc are installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "📦 Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "📦 Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Generate for auth-service
echo "📝 Generating auth-service protobuf..."
cd auth-service
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       internal/authpb/auth.proto

echo "✅ Generated:"
echo "  - internal/authpb/auth.pb.go"
echo "  - internal/authpb/auth_grpc.pb.go"
cd ..

# Generate for core-service
echo "📝 Generating core-service protobuf..."
cd core-service
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       internal/corepb/auth.proto

echo "✅ Generated:"
echo "  - internal/corepb/auth.pb.go"
echo "  - internal/corepb/auth_grpc.pb.go"
cd ..

echo "🎉 All protobuf files generated successfully!"