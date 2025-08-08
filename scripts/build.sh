#!/bin/bash

# Скрипт для сборки всех сервисов
set -e

echo "🔨 Building all services..."

# Генерация protobuf файлов
source scripts/generate-proto.sh

# Синхронизация зависимостей
echo "📦 Syncing dependencies..."
cd core-service && go mod tidy && cd ..
cd auth-service && go mod tidy && cd ..

# Сборка Docker образов
echo "🐳 Building Docker images..."
docker-compose build --parallel

echo "✅ All services built successfully!"