FROM golang:1.24-alpine

RUN apk add --no-cache git libc6-compat

WORKDIR /app

# Сначала копируем только файлы модулей
COPY go.mod go.sum ./

# Загружаем зависимости (этот слой закэшируется, если go.mod/go.sum не изменятся)
RUN go mod download

# Копируем остальные файлы
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 go build -o /app/subs-service /app/cmd/subs-service/main.go && \
    chmod +x /app/subs-service

EXPOSE 8080

CMD ["/app/subs-service"]