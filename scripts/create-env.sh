#!/bin/bash

# Script to create .env file
set -e

echo "🔧 Creating .env file..."
echo "=========================="

# Create .env file
cat > .env << 'EOF'
# Database Configuration
AUTH_DB_HOST=auth-db
AUTH_DB_PORT=5432
AUTH_DB_USER=auth_user
AUTH_DB_PASSWORD=auth_pass
AUTH_DB_NAME=auth_db
AUTH_DB_SSLMODE=disable

CORE_DB_HOST=core-db
CORE_DB_PORT=5432
CORE_DB_USER=core_user
CORE_DB_PASSWORD=core_pass
CORE_DB_NAME=core_db
CORE_DB_SSLMODE=disable

NOTIFY_DB_HOST=notify-db
NOTIFY_DB_PORT=5432
NOTIFY_DB_USER=notify_user
NOTIFY_DB_PASSWORD=notify_pass
NOTIFY_DB_NAME=notify_db
NOTIFY_DB_SSLMODE=disable

# Service Hosts (service names in Docker network)
AUTH_SERVICE_HOST=auth-service
CORE_SERVICE_HOST=core-service
NOTIFY_SERVICE_HOST=notification-service
RABBITMQ_HOST=rabbitmq

# Service Ports (single ports for internal and external access)
AUTH_SERVICE_PORT=50051
AUTH_HEALTH_PORT=8081
CORE_SERVICE_PORT=8080
NOTIFY_SERVICE_PORT=8082
RABBITMQ_PORT=5672
RABBITMQ_MANAGEMENT_PORT=15672

# Database Ports (internal ports in containers)
AUTH_DB_PORT=5432
CORE_DB_PORT=5432
NOTIFY_DB_PORT=5432

# Database External Ports (external ports for host access)
AUTH_DB_EXTERNAL_PORT=5433
CORE_DB_EXTERNAL_PORT=5434
NOTIFY_DB_EXTERNAL_PORT=5435

# Auth Service Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION=24h

# RabbitMQ Configuration
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_EXCHANGE=user_events
RABBITMQ_QUEUE=user_created

# TLS Configuration
ENABLE_TLS=false
TLS_CERT_FILE=/app/certs/server-cert.pem
TLS_KEY_FILE=/app/certs/server-key.pem

# Cookie Configuration
COOKIE_DOMAIN=localhost
COOKIE_MAX_AGE=3600

# Core Service Timeouts
CORE_READ_TIMEOUT=10s
CORE_WRITE_TIMEOUT=15s
CORE_IDLE_TIMEOUT=60s
CORE_SHUTDOWN_TIMEOUT=5s

# Core Service Rate Limiting
CORE_RATE_LIMIT_RPS=10
CORE_RATE_LIMIT_BURST=20

# Auth Service Timeouts
AUTH_SHUTDOWN_TIMEOUT=10s

# Notification Service Timeouts
NOTIFY_SHUTDOWN_TIMEOUT=10s
EOF

echo "✅ .env file created successfully!"

echo
echo "🔍 Validation of the created file..."
echo "================================"
./scripts/validate-env.sh

echo
echo "📋 File content:"
echo "===================="
cat .env

echo
echo "🚀 Now you can start the services:"
echo "   docker-compose up -d"
