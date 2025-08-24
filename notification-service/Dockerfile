FROM golang:1.24-alpine

# Install build dependencies
RUN apk add --no-cache git libc6-compat

WORKDIR /app
COPY notification-service/go.mod notification-service/go.sum ./notification-service/
WORKDIR /app/notification-service
RUN go mod download

WORKDIR /app
COPY notification-service ./notification-service

WORKDIR /app/notification-service

# Build the binary and make it executable
RUN CGO_ENABLED=0 go build -o notification-service ./cmd/notification-service && \
    chmod +x notification-service

# Create non-root user and set ownership
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup && \
    chown -R appuser:appgroup /app

# Ports are configured at runtime via compose/env; EXPOSE is static for metadata
EXPOSE 8082

# Switch to non-root user
USER appuser

CMD ["./notification-service"]
