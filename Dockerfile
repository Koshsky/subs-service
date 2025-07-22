FROM golang:1.24-alpine

RUN apk add --no-cache git libc6-compat

WORKDIR /app

COPY . .

RUN go mod tidy && \
    CGO_ENABLED=0 go build -o /app/subs-service /app/cmd/subs-service/main.go

RUN chmod +x /app/subs-service

EXPOSE 8080

CMD ["/app/subs-service"]