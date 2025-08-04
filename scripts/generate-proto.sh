#!/bin/bash

# Скрипт для генерации protobuf файлов
set -e

echo "🔧 Generating protobuf files..."

# Проверяем наличие protoc
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc не найден. Установите Protocol Buffers compiler:"
    echo "   Ubuntu/Debian: sudo apt install protobuf-compiler"
    echo "   macOS: brew install protobuf"
    exit 1
fi

# Проверяем наличие protoc-gen-go и protoc-gen-go-grpc
if ! command -v protoc-gen-go &> /dev/null; then
    echo "📦 Устанавливаем protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "📦 Устанавливаем protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Генерируем для auth-service
echo "📝 Generating auth-service protobuf..."
cd auth-service
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       internal/authpb/auth.proto

echo "✅ Generated:"
echo "  - internal/authpb/auth.pb.go"
echo "  - internal/authpb/auth_grpc.pb.go"
cd ..

# Генерируем для core-service
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