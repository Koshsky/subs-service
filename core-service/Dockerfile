FROM golang:1.24-alpine

RUN apk add --no-cache git libc6-compat wget

WORKDIR /app
COPY core-service/go.mod core-service/go.sum ./core-service/
WORKDIR /app/core-service
RUN go mod download

WORKDIR /app
COPY core-service ./core-service

# Create certs directory for TLS certificates
RUN mkdir -p /app/certs

WORKDIR /app/core-service

RUN CGO_ENABLED=0 go build -o core-service cmd/subs-service/main.go && \
    chmod +x core-service

# Create a non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# Change ownership of the application to the non-root user
RUN chown -R appuser:appgroup /app

ENV CORE_PORT=${CORE_PORT:-8080}
ENV DATABASE_URL=${DATABASE_URL:-postgres://user:password@localhost/subs_db?sslmode=disable}
ENV AUTH_SERVICE_ADDR=${AUTH_SERVICE_ADDR:-localhost:50051}
ENV ENABLE_TLS=${ENABLE_TLS:-false}
ENV TLS_CERT_FILE=${TLS_CERT_FILE:-/app/certs/server-cert.pem}

EXPOSE ${CORE_PORT}

# Switch to non-root user
USER appuser

CMD ["./core-service"]