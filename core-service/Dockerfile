FROM golang:1.24-alpine

# Install build dependencies and prepare certs directory
RUN apk add --no-cache git libc6-compat wget && \
    mkdir -p /app/certs

WORKDIR /app
COPY core-service/go.mod core-service/go.sum ./core-service/
WORKDIR /app/core-service
RUN go mod download

WORKDIR /app
COPY core-service ./core-service
# Copy TLS certs from repo (if present)
COPY certs /app/certs

WORKDIR /app/core-service

# Build the binary and make it executable
RUN CGO_ENABLED=0 go build -o core-service ./cmd/core-service && \
    chmod +x core-service

# Create non-root user and set ownership
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup && \
    chown -R appuser:appgroup /app

# Ports are configured at runtime via compose/env; EXPOSE is static for metadata
EXPOSE 8080

# Switch to non-root user
USER appuser

CMD ["./core-service"]