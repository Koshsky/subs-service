#!/bin/bash

# Script to validate critical environment variables
set -e

echo "🔍 Validating Critical Environment Variables"
echo "============================================="

# Array of critical environment variables that MUST be set
CRITICAL_VARS=(
    # Database Configuration
    "AUTH_DB_USER"
    "AUTH_DB_PASSWORD"
    "AUTH_DB_NAME"
    "CORE_DB_USER"
    "CORE_DB_PASSWORD"
    "CORE_DB_NAME"
    "NOTIFY_DB_USER"
    "NOTIFY_DB_PASSWORD"
    "NOTIFY_DB_NAME"

    # Service Ports
    "AUTH_SERVICE_PORT"
    "CORE_SERVICE_PORT"
    "NOTIFY_SERVICE_PORT"
    "RABBITMQ_PORT"
    "RABBITMQ_MANAGEMENT_PORT"

    # Security
    "JWT_SECRET"
    "ENABLE_TLS"

    # RabbitMQ Configuration
    "RABBITMQ_USER"
    "RABBITMQ_PASSWORD"
    "RABBITMQ_EXCHANGE"
    "RABBITMQ_QUEUE"

    # Cookie Configuration
    "COOKIE_DOMAIN"
    "COOKIE_MAX_AGE"
)

# Array of variables that should have specific values or formats
VALIDATION_RULES=(
    "JWT_SECRET:min_length:32"
    "AUTH_SERVICE_PORT:port"
    "CORE_SERVICE_PORT:port"
    "NOTIFY_SERVICE_PORT:port"
    "RABBITMQ_PORT:port"
    "RABBITMQ_MANAGEMENT_PORT:port"
    "ENABLE_TLS:boolean"
    "COOKIE_MAX_AGE:positive_integer"
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

# Function to validate positive integer
validate_positive_integer() {
    local value=$1
    local var_name=$2

    if [[ "$value" =~ ^[1-9][0-9]*$ ]]; then
        echo "✅ $var_name positive integer $value is valid"
        return 0
    else
        echo "❌ ERROR: $var_name value $value is not a positive integer"
        return 1
    fi
}

# Function to check for default values in code
check_default_values() {
    echo
    echo "🔍 Checking for default values in code..."
    echo "=========================================="

    local found_defaults=false

    # Check for common default value patterns
    local patterns=(
        ":-default"
        ":-guest"
        ":-auth_user"
        ":-core_user"
        ":-notify_user"
        ":-auth_pass"
        ":-core_pass"
        ":-notify_pass"
        ":-auth_db"
        ":-core_db"
        ":-notify_db"
        ":-50051"
        ":-8080"
        ":-8081"
        ":-8082"
        ":-5672"
        ":-15672"
        ":-5433"
        ":-5434"
        ":-5435"
        ":-3600"
        ":-localhost"
    )

    for pattern in "${patterns[@]}"; do
        if grep -r "$pattern" . --include="*.go" --include="*.yaml" --include="*.yml" --include="*.sh" --exclude-dir=.git --exclude-dir=tmp > /dev/null 2>&1; then
            echo "⚠️  WARNING: Found default value pattern '$pattern' in code"
            found_defaults=true
        fi
    done

    if [ "$found_defaults" = false ]; then
        echo "✅ No default value patterns found in code"
    fi
}

# Function to check for hardcoded values
check_hardcoded_values() {
    echo
    echo "🔍 Checking for hardcoded critical values..."
    echo "============================================"

    local found_hardcoded=false

    # Check for hardcoded critical values
    local hardcoded_patterns=(
        "guest:guest"
        "auth_user"
        "core_user"
        "notify_user"
        "auth_pass"
        "core_pass"
        "notify_pass"
        "auth_db"
        "core_db"
        "notify_db"
        "50051"
        "8080"
        "8081"
        "8082"
        "5672"
        "15672"
        "5433"
        "5434"
        "5435"
        "3600"
        "localhost"
    )

    for pattern in "${hardcoded_patterns[@]}"; do
        if grep -r "$pattern" . --include="*.go" --include="*.yaml" --include="*.yml" --include="*.sh" --exclude-dir=.git --exclude-dir=tmp > /dev/null 2>&1; then
            echo "⚠️  WARNING: Found hardcoded value '$pattern' in code"
            found_hardcoded=true
        fi
    done

    if [ "$found_hardcoded" = false ]; then
        echo "✅ No hardcoded critical values found in code"
    fi
}

# Main validation logic
main() {
    local exit_code=0
    local errors_found=false

    # Check if .env file exists
    if [ ! -f .env ]; then
        echo "❌ CRITICAL ERROR: .env file not found"
        echo "   Run: ./scripts/create-env.sh"
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
            "positive_integer")
                if ! validate_positive_integer "$var_value" "$var_name"; then
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

    # Check for port conflicts
    local ports=(
        "$AUTH_SERVICE_PORT"
        "$CORE_SERVICE_PORT"
        "$NOTIFY_SERVICE_PORT"
        "$RABBITMQ_PORT"
        "$RABBITMQ_MANAGEMENT_PORT"
    )

    local unique_ports=($(printf '%s\n' "${ports[@]}" | sort -u))
    local total_ports=${#ports[@]}
    local unique_count=${#unique_ports[@]}

    if [ "$total_ports" -eq "$unique_count" ]; then
        echo "✅ No port conflicts detected"
    else
        echo "❌ ERROR: Port conflicts detected"
        errors_found=true
        exit_code=1
fi

echo
    echo "🔍 Security checks..."
    echo "===================="

    # Security checks
if [ "$JWT_SECRET" = "your-super-secret-jwt-key-change-in-production" ]; then
        echo "⚠️  WARNING: JWT_SECRET is using default value (OK for development)"
    else
        echo "✅ JWT_SECRET is properly configured"
    fi

    if [ "$AUTH_DB_PASSWORD" = "auth_pass" ] || [ "$CORE_DB_PASSWORD" = "core_pass" ] || [ "$NOTIFY_DB_PASSWORD" = "notify_pass" ]; then
        echo "⚠️  WARNING: Database passwords are using default values (OK for development)"
    else
        echo "✅ Database passwords are properly configured"
    fi

    if [ "$RABBITMQ_PASSWORD" = "guest" ]; then
        echo "⚠️  WARNING: RabbitMQ password is using default value (OK for development)"
    else
        echo "✅ RabbitMQ password is properly configured"
    fi

    echo
    echo "📊 Validation Summary"
    echo "===================="

    if [ "$errors_found" = true ]; then
        echo "❌ Validation failed with errors"
        echo
        echo "💡 To fix issues:"
        echo "   1. Update .env file with proper values"
        echo "   2. Remove default values from code"
        echo "   3. Replace hardcoded values with environment variables"
        echo "   4. Run validation again: ./scripts/validate-env.sh"
    else
        echo "✅ All validations passed successfully!"
        echo
        echo "🚀 Environment is ready for deployment"
    fi

    exit $exit_code
}

# Run main function
main "$@"
