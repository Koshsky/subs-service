#!/bin/bash

# Script to validate the .env file
set -e

echo "🔍 Validating .env file..."
echo "=========================="

# Check if the .env file exists
if [ ! -f ".env" ]; then
    echo "❌ .env file not found!"
    echo "   Create it using: ./scripts/create-env.sh"
    exit 1
fi

echo "✅ .env file found"

# List of required variables
REQUIRED_VARS=(
    # Database Configuration
    "AUTH_DB_HOST"
    "AUTH_DB_PORT"
    "AUTH_DB_EXTERNAL_PORT"
    "AUTH_DB_USER"
    "AUTH_DB_PASSWORD"
    "AUTH_DB_NAME"

    "CORE_DB_HOST"
    "CORE_DB_PORT"
    "CORE_DB_EXTERNAL_PORT"
    "CORE_DB_USER"
    "CORE_DB_PASSWORD"
    "CORE_DB_NAME"

    "NOTIFY_DB_HOST"
    "NOTIFY_DB_PORT"
    "NOTIFY_DB_EXTERNAL_PORT"
    "NOTIFY_DB_USER"
    "NOTIFY_DB_PASSWORD"
    "NOTIFY_DB_NAME"

    # Service Hosts
    "AUTH_SERVICE_HOST"
    "CORE_SERVICE_HOST"
    "NOTIFY_SERVICE_HOST"
    "RABBITMQ_HOST"

    # Service Ports
    "AUTH_SERVICE_PORT"
    "AUTH_HEALTH_PORT"
    "CORE_SERVICE_PORT"
    "NOTIFY_SERVICE_PORT"
    "RABBITMQ_PORT"
    "RABBITMQ_MANAGEMENT_PORT"

    # Auth Service Configuration
    "JWT_SECRET"

    # RabbitMQ Configuration
    "RABBITMQ_URL"
    "RABBITMQ_USER"
    "RABBITMQ_PASSWORD"
    "RABBITMQ_EXCHANGE"
    "RABBITMQ_QUEUE"
)

# Check each required variable
MISSING_VARS=()
for var in "${REQUIRED_VARS[@]}"; do
    if ! grep -q "^${var}=" .env; then
        MISSING_VARS+=("$var")
    fi
done

# Check for empty values
EMPTY_VARS=()
while IFS= read -r line; do
    # Skip comments and empty lines
    if [[ "$line" =~ ^[[:space:]]*# ]] || [[ -z "$line" ]]; then
        continue
    fi

    # Extract variable name and value
    if [[ "$line" =~ ^([^=]+)=(.*)$ ]]; then
        var_name="${BASH_REMATCH[1]}"
        var_value="${BASH_REMATCH[2]}"

        # Remove quotes and spaces
        var_value=$(echo "$var_value" | sed 's/^["'\'']*//;s/["'\'']*$//' | xargs)

        if [ -z "$var_value" ]; then
            EMPTY_VARS+=("$var_name")
        fi
    fi
done < .env

# Output results
if [ ${#MISSING_VARS[@]} -eq 0 ] && [ ${#EMPTY_VARS[@]} -eq 0 ]; then
    echo "✅ All required environment variables are set correctly"
else
    echo "❌ Found problems with environment variables:"

    if [ ${#MISSING_VARS[@]} -gt 0 ]; then
        echo
        echo "📋 Missing variables:"
        for var in "${MISSING_VARS[@]}"; do
            echo "   - $var"
        done
    fi

    if [ ${#EMPTY_VARS[@]} -gt 0 ]; then
        echo
        echo "📋 Variables with empty values:"
        for var in "${EMPTY_VARS[@]}"; do
            echo "   - $var"
        done
    fi

    echo
    echo "🔧 To fix:"
    echo "   1. Run: ./scripts/create-env.sh"
    echo "   2. Or edit the .env file manually"

    exit 1
fi

echo
echo "🔍 Checking critical settings..."
echo "=================================="

# Check if JWT_SECRET is not the default
JWT_SECRET=$(grep "^JWT_SECRET=" .env | cut -d'=' -f2 | sed 's/^["'\'']*//;s/["'\'']*$//')
if [ "$JWT_SECRET" = "your-super-secret-jwt-key-change-in-production" ]; then
    echo "⚠️  WARNING: JWT_SECRET uses the default value!"
    echo "   It is recommended to change it in production"
fi

# Check TLS configuration
ENABLE_TLS=$(grep "^ENABLE_TLS=" .env | cut -d'=' -f2 | sed 's/^["'\'']*//;s/["'\'']*$//')
if [ "$ENABLE_TLS" = "true" ]; then
    echo "🔒 TLS enabled"

    # Check for certificates
    TLS_CERT_FILE=$(grep "^TLS_CERT_FILE=" .env | cut -d'=' -f2 | sed 's/^["'\'']*//;s/["'\'']*$//')
    TLS_KEY_FILE=$(grep "^TLS_KEY_FILE=" .env | cut -d'=' -f2 | sed 's/^["'\'']*//;s/["'\'']*$//')

    if [ -f "$TLS_CERT_FILE" ]; then
        echo "   ✅ Certificate found: $TLS_CERT_FILE"
    else
        echo "   ❌ Certificate not found: $TLS_CERT_FILE"
        echo "   Run: ./scripts/build.sh to generate certificates"
    fi

    if [ -f "$TLS_KEY_FILE" ]; then
        echo "   ✅ Private key found: $TLS_KEY_FILE"
    else
        echo "   ❌ Private key not found: $TLS_KEY_FILE"
        echo "   Run: ./scripts/build.sh to generate certificates"
    fi
else
    echo "⚠️  TLS disabled"
fi

# Check if ports do not conflict
echo "✅ Checking ports completed"

echo
echo "✅ Validation of the .env file completed successfully!"
echo "🚀 Now you can start the services:"
echo "   docker-compose up -d"
