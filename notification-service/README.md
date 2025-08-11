# Notification Service

Сервис уведомлений для обработки событий из RabbitMQ и отправки уведомлений пользователям.

## Функциональность

- Обработка событий `user.created` из RabbitMQ
- Логирование событий создания пользователей
- Готовность к расширению для отправки email/SMS уведомлений

## Архитектура

### RabbitMQ
- Exchange: `user_events` (topic)
- Queue: `user_created`
- Routing Key: `user.created`

### События
```json
{
  "user_id": "uuid",
  "email": "user@example.com"
}
```

## Конфигурация

### Переменные окружения
- `NOTIFY_DB_HOST` - хост базы данных (по умолчанию: notify-db)
- `NOTIFY_DB_PORT` - порт базы данных (по умолчанию: 5432)
- `NOTIFY_DB_USER` - пользователь БД (по умолчанию: notify_user)
- `NOTIFY_DB_PASSWORD` - пароль БД (по умолчанию: notify_pass)
- `NOTIFY_DB_NAME` - имя БД (по умолчанию: notify_db)
- `NOTIFY_PORT` - порт HTTP сервера (по умолчанию: 8082)
- `RABBITMQ_URL` - URL RabbitMQ (по умолчанию: amqp://guest:guest@rabbitmq:5672/)
- `RABBITMQ_EXCHANGE` - exchange RabbitMQ (по умолчанию: user_events)
- `RABBITMQ_QUEUE` - queue RabbitMQ (по умолчанию: user_created)

## Запуск

```bash
# Локально
go run cmd/notification-service/main.go

# Docker
docker-compose up notification-service
```

## Health Check

```
GET /health
```

Ответ:
```json
{
  "status": "ok",
  "service": "notification-service",
  "timestamp": "2024-01-01T12:00:00Z"
}
```
