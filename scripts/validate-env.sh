#!/bin/bash

# Script to validate critical environment variables
set -e

echo "🔍 Validating Critical Environment Variables"
echo "============================================="

# Array of critical environment variables that MUST be set
CRITICAL_VARS=(
    # Database Configuration (REQUIRED by code)
    "AUTH_DB_USER"
    "AUTH_DB_PASSWORD"
    "AUTH_DB_NAME"
    "AUTH_DB_PORT"
    "CORE_DB_USER"
    "CORE_DB_PASSWORD"
    "CORE_DB_NAME"
    "CORE_DB_PORT"
    "NOTIFY_DB_USER"
    "NOTIFY_DB_PASSWORD"
    "NOTIFY_DB_NAME"
    "NOTIFY_DB_PORT"

    # Service Ports (REQUIRED by code)
    "AUTH_SERVICE_PORT"
    "CORE_SERVICE_PORT"
    "NOTIFY_SERVICE_PORT"
    "RABBITMQ_PORT"
    "RABBITMQ_MANAGEMENT_PORT"

    # Security (REQUIRED by code)
    "JWT_SECRET"
    "ENABLE_TLS"

    # RabbitMQ Configuration
    "RABBITMQ_USER"
    "RABBITMQ_PASSWORD"
    "RABBITMQ_EXCHANGE"
    "RABBITMQ_QUEUE"
)

# Array of variables that should have specific values or formats
VALIDATION_RULES=(
    "JWT_SECRET:min_length:32"
    "AUTH_SERVICE_PORT:port"
    "CORE_SERVICE_PORT:port"
    "AUTH_DB_PORT:port"
    "CORE_DB_PORT:port"
    "NOTIFY_DB_PORT:port"
    "NOTIFY_SERVICE_PORT:port"
    "ENABLE_TLS:boolean"
)

# Function to check if variable is set
check_variable_set() {
    local var_name=$1
    local var_value="${!var_name}"

    if [ -z "$var_value" ]; then
        echo "❌ CRITICAL ERROR: $var_name is not set"
        return 1
    else
        echo "✅ $var_name is set"
        return 0
    fi
}

# Function to validate port numbers
validate_port() {
    local port=$1
    local var_name=$2

    if [[ "$port" =~ ^[0-9]+$ ]] && [ "$port" -ge 1024 ] && [ "$port" -le 65535 ]; then
        echo "✅ $var_name port $port is valid"
        return 0
    else
        echo "❌ ERROR: $var_name port $port is invalid (must be 1024-65535)"
        return 1
    fi
}

# Function to validate boolean values
validate_boolean() {
    local value=$1
    local var_name=$2

    if [[ "$value" =~ ^(true|false)$ ]]; then
        echo "✅ $var_name boolean value $value is valid"
        return 0
    else
        echo "❌ ERROR: $var_name boolean value $value is invalid (must be true or false)"
        return 1
    fi
}

# Function to validate minimum length
validate_min_length() {
    local value=$1
    local min_length=$2
    local var_name=$3

    if [ ${#value} -ge "$min_length" ]; then
        echo "✅ $var_name meets minimum length requirement ($min_length chars)"
        return 0
    else
        echo "❌ ERROR: $var_name is too short (${#value} chars, minimum $min_length required)"
        return 1
    fi
}

# Function to check for default values in code
check_default_values() {
    echo
    echo "🔍 Checking for default values in code..."
    echo "========================================="

    # Check for default values that should be overridden in production
    local default_vars=(
        "AUTH_DB_HOST:auth-db"
        "CORE_DB_HOST:core-db"
        "NOTIFY_DB_HOST:notify-db"
        "AUTH_DB_SSLMODE:disable"
        "CORE_DB_SSLMODE:disable"
        "NOTIFY_DB_SSLMODE:disable"
        "NOTIFY_SERVICE_PORT:8082"
        "RABBITMQ_URL:amqp://guest:guest@rabbitmq:5672/"
        "NOTIFY_SHUTDOWN_TIMEOUT:10s"
        "TLS_CERT_FILE:certs/server-cert.pem"
        "TLS_KEY_FILE:certs/server-key.pem"
    )

    for default_var in "${default_vars[@]}"; do
        IFS=':' read -r var_name default_value <<< "$default_var"
        local var_value="${!var_name}"

        if [ "$var_value" = "$default_value" ]; then
            echo "⚠️  WARNING: $var_name is using default value '$default_value'"
        else
            echo "✅ $var_name is customized"
        fi
    done
}

# Function to check for hardcoded values
check_hardcoded_values() {
    echo
    echo "🔍 Checking for hardcoded values..."
    echo "=================================="

    # Check for development passwords
    local weak_passwords=(
        "auth_pass"
        "core_pass"
        "notify_pass"
        "guest"
    )

    for password in "${weak_passwords[@]}"; do
        if [[ "$AUTH_DB_PASSWORD" == "$password" ]] || \
           [[ "$CORE_DB_PASSWORD" == "$password" ]] || \
           [[ "$NOTIFY_DB_PASSWORD" == "$password" ]]; then
            echo "⚠️  WARNING: Using weak password '$password' - change in production"
        fi
    done

    # Check for weak JWT secret
    if [[ "$JWT_SECRET" == "your-super-secret-jwt-key-change-in-production"* ]]; then
        echo "⚠️  WARNING: Using default JWT secret - change in production"
    fi
}

# Main validation logic
main() {
    local exit_code=0
    local errors_found=false

    # Check if .env file exists
    if [ ! -f .env ]; then
        echo "❌ CRITICAL ERROR: .env file not found"
        echo "   Run: ./scripts/generate-env.sh"
        exit 1
    fi

    # Source the .env file
    set -a
    source .env
    set +a

    echo "📋 Checking critical variables..."
    echo "================================"

    # Check all critical variables are set
    for var in "${CRITICAL_VARS[@]}"; do
        if ! check_variable_set "$var"; then
            errors_found=true
            exit_code=1
        fi
    done

    echo
    echo "🔧 Validating variable formats..."
    echo "================================"

    # Apply validation rules
    for rule in "${VALIDATION_RULES[@]}"; do
        IFS=':' read -r var_name rule_type rule_value <<< "$rule"
        local var_value="${!var_name}"

        case $rule_type in
            "port")
                if ! validate_port "$var_value" "$var_name"; then
                    errors_found=true
                    exit_code=1
                fi
                ;;
            "boolean")
                if ! validate_boolean "$var_value" "$var_name"; then
                    errors_found=true
                    exit_code=1
                fi
                ;;
            "min_length")
                if ! validate_min_length "$var_value" "$rule_value" "$var_name"; then
                    errors_found=true
                    exit_code=1
                fi
                ;;
        esac
    done

    # Check for default values and hardcoded values
    check_default_values
    check_hardcoded_values

    echo
    echo "🔍 Checking port conflicts..."
    echo "============================"

    # Check for port conflicts (only service ports, not database ports)
    local service_ports=(
        "$AUTH_SERVICE_PORT"
        "$CORE_SERVICE_PORT"
        "$NOTIFY_SERVICE_PORT"
        "$RABBITMQ_PORT"
        "$RABBITMQ_MANAGEMENT_PORT"
    )

    local unique_ports=()
    for port in "${service_ports[@]}"; do
        if [[ " ${unique_ports[@]} " =~ " ${port} " ]]; then
            echo "❌ ERROR: Port conflict detected - port $port is used multiple times"
            errors_found=true
            exit_code=1
        else
            unique_ports+=("$port")
            echo "✅ Port $port is unique"
        fi
    done

    # Note about database ports
    echo "ℹ️  Database ports (AUTH_DB_PORT, CORE_DB_PORT, NOTIFY_DB_PORT) can be the same as they run in separate containers"

    echo
    echo "📊 Validation Summary"
    echo "===================="

    if [ "$errors_found" = true ]; then
        echo "❌ Validation failed - please fix the errors above"
        echo "💡 Run ./scripts/generate-env.sh to regenerate with correct values"
        exit $exit_code
    else
        echo "✅ All validations passed!"
        echo "🚀 Environment is ready for use"
    fi
}

# Run main function
main "$@"
