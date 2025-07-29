FROM golang:1.24-alpine

RUN apk add --no-cache git libc6-compat

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/subs-service /app/cmd/subs-service/main.go && \
    chmod +x /app/subs-service

COPY .env ./

EXPOSE 8080

CMD ["/app/subs-service"]