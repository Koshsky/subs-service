#!/bin/bash

# Script to create .env file
set -e

echo "🔧 Creating .env file..."
echo "=========================="

# Create .env file
cat > .env << 'EOF'
# =============================================================================
# SUBS-SERVICE ENVIRONMENT CONFIGURATION
# =============================================================================
# Generated on: $(date)
# WARNING: Change these values in production!
# =============================================================================

# =============================================================================
# REQUIRED VARIABLES (will cause panic if not set)
# =============================================================================

# Database Configuration - REQUIRED
AUTH_DB_USER=auth_user
AUTH_DB_PASSWORD=auth_pass
AUTH_DB_NAME=auth_db
AUTH_DB_PORT=5432

CORE_DB_USER=core_user
CORE_DB_PASSWORD=core_pass
CORE_DB_NAME=core_db
CORE_DB_PORT=5432

NOTIFY_DB_USER=notify_user
NOTIFY_DB_PASSWORD=notify_pass
NOTIFY_DB_NAME=notify_db
NOTIFY_DB_PORT=5432

# Service Ports - REQUIRED
AUTH_SERVICE_PORT=50051
CORE_SERVICE_PORT=8080
NOTIFY_SERVICE_PORT=8082

# Security - REQUIRED (minimum 32 characters)
JWT_SECRET=your-super-secret-jwt-key-change-in-production-minimum-32-chars
ENABLE_TLS=false

# RabbitMQ Configuration - REQUIRED
RABBITMQ_PORT=5672
RABBITMQ_MANAGEMENT_PORT=15672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_EXCHANGE=user_events
RABBITMQ_QUEUE=user_created

# =============================================================================
# OPTIONAL VARIABLES (have default values)
# =============================================================================

# Database Hosts (optional - have defaults)
AUTH_DB_HOST=auth-db
CORE_DB_HOST=core-db
NOTIFY_DB_HOST=notify-db

# Database SSL Mode (optional - have defaults)
AUTH_DB_SSLMODE=disable
CORE_DB_SSLMODE=disable
NOTIFY_DB_SSLMODE=disable

# Notification Service Configuration (optional - has default)
NOTIFY_SHUTDOWN_TIMEOUT=10s

# RabbitMQ Configuration (optional - have defaults)
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

# TLS Configuration (optional - have defaults)
TLS_CERT_FILE=certs/server-cert.pem
TLS_KEY_FILE=certs/server-key.pem

# =============================================================================
# DOCKER-COMPOSE ONLY VARIABLES (not used in Go code)
# =============================================================================

# Cookie Configuration (used only in docker-compose.yaml)
COOKIE_DOMAIN=
COOKIE_MAX_AGE=

# JWT Expiration (used only in docker-compose.yaml)
JWT_EXPIRATION=

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

echo
echo "⚠️  IMPORTANT SECURITY NOTES:"
echo "   1. This .env file contains DEFAULT VALUES for development"
echo "   2. Change all passwords and secrets in production"
echo "   3. Use strong JWT secrets (minimum 32 characters)"
echo "   4. Enable TLS in production environments"
echo "   5. Never commit .env files to version control"
echo "   6. Run ./scripts/validate-env.sh to check for issues"
echo
echo "📚 For more information, see:"
echo "   - docs/ENVIRONMENT.md"
echo "   - docs/RABBITMQ_BEHAVIOR.md"
echo "   - docs/SECURITY.md"
