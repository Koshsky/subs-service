# Environment Variables Documentation

## Overview

This document describes the environment variables used in the subscription service and their critical importance for security and functionality.

## Critical Environment Variables

These variables **MUST** be set and **CANNOT** have default values in production:

### Security Variables

| Variable | Description | Validation | Example |
|----------|-------------|------------|---------|
| `JWT_SECRET` | Secret key for JWT token signing | Min 32 characters | `your-super-secret-jwt-key-change-in-production` |
| `AUTH_DB_PASSWORD` | Auth database password | Non-empty | `secure_auth_password_123` |
| `CORE_DB_PASSWORD` | Core database password | Non-empty | `secure_core_password_123` |
| `NOTIFY_DB_PASSWORD` | Notification database password | Non-empty | `secure_notify_password_123` |
| `RABBITMQ_PASSWORD` | RabbitMQ password | Non-empty | `secure_rabbitmq_password_123` |

### Port Configuration

| Variable | Description | Validation | Example |
|----------|-------------|------------|---------|
| `AUTH_SERVICE_PORT` | Auth service gRPC port | 1024-65535 | `50051` |
| `CORE_SERVICE_PORT` | Core service HTTP port | 1024-65535 | `8080` |
| `NOTIFY_SERVICE_PORT` | Notification service port | 1024-65535 | `8082` |
| `RABBITMQ_PORT` | RabbitMQ AMQP port | 1024-65535 | `5672` |
| `RABBITMQ_MANAGEMENT_PORT` | RabbitMQ management port | 1024-65535 | `15672` |

### Database Configuration

| Variable | Description | Validation | Example |
|----------|-------------|------------|---------|
| `AUTH_DB_NAME` | Auth database name | Non-empty | `auth_db` |
| `CORE_DB_NAME` | Core database name | Non-empty | `core_db` |
| `NOTIFY_DB_NAME` | Notification database name | Non-empty | `notify_db` |
| `AUTH_DB_USER` | Auth database user | Non-empty | `auth_user` |
| `CORE_DB_USER` | Core database user | Non-empty | `core_user` |
| `NOTIFY_DB_USER` | Notification database user | Non-empty | `notify_user` |

### RabbitMQ Configuration

| Variable | Description | Validation | Example |
|----------|-------------|------------|---------|
| `RABBITMQ_USER` | RabbitMQ username | Non-empty | `rabbitmq_user` |
| `RABBITMQ_EXCHANGE` | RabbitMQ exchange name | Non-empty | `user_events` |
| `RABBITMQ_QUEUE` | RabbitMQ queue name | Non-empty | `user_created` |

### TLS Configuration

| Variable | Description | Validation | Example |
|----------|-------------|------------|---------|
| `ENABLE_TLS` | Enable TLS encryption | Boolean | `false` |

## Non-Critical Environment Variables

These variables can have default values:

### Database Configuration (Internal)

| Variable | Description | Default |
|----------|-------------|---------|
| `AUTH_DB_HOST` | Auth database host | `auth-db` |
| `AUTH_DB_PORT` | Auth database internal port | `5432` |
| `AUTH_DB_SSLMODE` | Auth database SSL mode | `disable` |
| `CORE_DB_HOST` | Core database host | `core-db` |
| `CORE_DB_PORT` | Core database internal port | `5432` |
| `CORE_DB_SSLMODE` | Core database SSL mode | `disable` |
| `NOTIFY_DB_HOST` | Notification database host | `notify-db` |
| `NOTIFY_DB_PORT` | Notification database internal port | `5432` |
| `NOTIFY_DB_SSLMODE` | Notification database SSL mode | `disable` |

### Service Hosts

| Variable | Description | Default |
|----------|-------------|---------|
| `RABBITMQ_HOST` | RabbitMQ host | `rabbitmq` |

### TLS Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `TLS_CERT_FILE` | TLS certificate file path | `certs/server-cert.pem` |
| `TLS_KEY_FILE` | TLS private key file path | `certs/server-key.pem` |

## Environment Setup

### Development Setup

For development, you can create a `.env` file with default values:

```bash
./scripts/generate-env.sh
```

This will create a `.env` file with development-friendly default values that you can then customize.

### Production Setup

For production, you **MUST**:

1. **Never use default values** for critical variables
2. **Use strong, unique passwords** for all databases and services
3. **Use a strong JWT secret** (minimum 32 characters)
4. **Enable TLS** for secure communication
5. **Change all default credentials**

Example production `.env` file:

```bash
# Security Variables
JWT_SECRET=your-super-secret-jwt-key-change-in-production-minimum-32-chars
AUTH_DB_PASSWORD=secure_auth_password_123
CORE_DB_PASSWORD=secure_core_password_123
NOTIFY_DB_PASSWORD=secure_notify_password_123
RABBITMQ_PASSWORD=secure_rabbitmq_password_123

# Port Configuration
AUTH_SERVICE_PORT=50051
CORE_SERVICE_PORT=8080
NOTIFY_SERVICE_PORT=8082
RABBITMQ_PORT=5672
RABBITMQ_MANAGEMENT_PORT=15672

# Database Configuration
AUTH_DB_NAME=auth_db
CORE_DB_NAME=core_db
NOTIFY_DB_NAME=notify_db
AUTH_DB_USER=auth_user
CORE_DB_USER=core_user
NOTIFY_DB_USER=notify_user

# RabbitMQ Configuration
RABBITMQ_USER=rabbitmq_user
RABBITMQ_EXCHANGE=user_events
RABBITMQ_QUEUE=user_created

# TLS Configuration
ENABLE_TLS=true
```

## Validation

The system includes comprehensive validation for environment variables:

### Build-Time Validation

You can validate environment variables manually:

```bash
./scripts/validate-env.sh
```

This will:
- Check that all critical variables are set
- Validate port numbers (1024-65535)
- Validate JWT secret length (minimum 32 characters)
- Check for port conflicts
- Warn about default values in development

### Manual Validation

You can manually validate environment variables:

```bash
./scripts/validate-env.sh
```

### Runtime Validation

The application will panic if critical environment variables are not set:

```go
// This will panic if AUTH_DB_PASSWORD is not set
password := utils.GetEnvRequired("AUTH_DB_PASSWORD")

// This will panic if JWT_SECRET is less than 32 characters
secret := utils.GetEnvRequiredWithValidation("JWT_SECRET", utils.ValidateMinLength(32))
```

## Security Best Practices

1. **Never commit `.env` files** to version control
2. **Use different passwords** for each database and service
3. **Rotate secrets regularly** in production
4. **Use environment-specific configurations**
5. **Enable TLS** in production environments
6. **Use strong, random passwords** (consider using password generators)
7. **Limit access** to environment variables to authorized personnel only

## Troubleshooting

### Common Issues

1. **"CRITICAL ERROR: Environment variable X is not set"**
   - Solution: Set the missing environment variable in your `.env` file

2. **"Port conflicts detected"**
   - Solution: Ensure all ports are unique across services

3. **"JWT_SECRET validation failed"**
   - Solution: Use a JWT secret with at least 32 characters

4. **"Database connection failed"**
   - Solution: Check database credentials and ensure database is running

### Validation Commands

```bash
# Check environment variables
./scripts/validate-env.sh

# Check service health
./scripts/health-check.sh

# Test API functionality
./scripts/test-api.sh
```

## Related Documentation

- [README.md](../README.md) - Project overview
- [SECURITY.md](../docs/SECURITY.md) - Security guidelines
- [TLS-SETUP.md](../docs/TLS-SETUP.md) - TLS configuration
- [HEALTHCHECK.md](../HEALTHCHECK.md) - Health check documentation
