FROM golang:1.24-alpine

RUN apk add --no-cache git libc6-compat

WORKDIR /app
COPY shared/go.mod shared/go.sum ./shared/
WORKDIR /app/shared
RUN go mod download

WORKDIR /app
COPY core-service/go.mod core-service/go.sum ./core-service/
WORKDIR /app/core-service
RUN go mod download

WORKDIR /app
COPY core-service ./core-service
COPY shared ./shared

WORKDIR /app/core-service

RUN CGO_ENABLED=0 go build -o core-service cmd/subs-service/main.go && \
    chmod +x core-service

ENV CORE_PORT=${CORE_PORT:-8080}
ENV DATABASE_URL=${DATABASE_URL:-postgres://user:password@localhost/subs_db?sslmode=disable}
ENV AUTH_SERVICE_ADDR=${AUTH_SERVICE_ADDR:-localhost:50051}

EXPOSE ${CORE_PORT}

CMD ["./core-service"]