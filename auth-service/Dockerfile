FROM golang:1.24-alpine

RUN apk add --no-cache git libc6-compat wget

WORKDIR /app
COPY shared/go.mod shared/go.sum ./shared/
WORKDIR /app/shared
RUN go mod download

WORKDIR /app
COPY auth-service/go.mod auth-service/go.sum ./auth-service/
WORKDIR /app/auth-service
RUN go mod download

WORKDIR /app
COPY auth-service ./auth-service
COPY shared ./shared

# Create certs directory for TLS certificates
RUN mkdir -p /app/certs

WORKDIR /app/auth-service

RUN CGO_ENABLED=0 go build -o auth-service cmd/auth-service/main.go && \
    chmod +x auth-service

# Create a non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# Change ownership of the application to the non-root user
RUN chown -R appuser:appgroup /app

ENV AUTH_PORT=${AUTH_PORT:-50051}
ENV DATABASE_URL=${DATABASE_URL:-postgres://user:password@localhost/subs_db?sslmode=disable}
ENV JWT_SECRET=${JWT_SECRET:-your_jwt_secret_key}
ENV ENABLE_TLS=${ENABLE_TLS:-false}
ENV TLS_CERT_FILE=${TLS_CERT_FILE:-/app/certs/server-cert.pem}
ENV TLS_KEY_FILE=${TLS_KEY_FILE:-/app/certs/server-key.pem}

EXPOSE ${AUTH_PORT}

# Switch to non-root user
USER appuser

CMD ["./auth-service"]