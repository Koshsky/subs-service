# Notification Service

Notification service for processing events from RabbitMQ and sending notifications to users.

## Functionality

- Processing `user.created` events from RabbitMQ
- Logging user creation events
- Ready for extension to send email/SMS notifications

## Architecture

### RabbitMQ
- Exchange: `user_events` (topic)
- Queue: `user_created`
- Routing Key: `user.created`

### Events
```json
{
  "user_id": "uuid",
  "email": "user@example.com"
}
```
jit warmup
## Configuration

### Environment Variables
- `NOTIFY_DB_HOST` - database host (default: notify-db)
- `NOTIFY_DB_PORT` - database port (default: 5432)
- `NOTIFY_DB_USER` - database user (default: notify_user)
- `NOTIFY_DB_PASSWORD` - database password (default: notify_pass)
- `NOTIFY_DB_NAME` - database name (default: notify_db)
- `NOTIFY_SERVICE_PORT` - HTTP server port (default: 8082)
- `RABBITMQ_URL` - RabbitMQ URL (default: amqp://guest:guest@rabbitmq:5672/)
- `RABBITMQ_EXCHANGE` - RabbitMQ exchange (default: user_events)
- `RABBITMQ_QUEUE` - RabbitMQ queue (default: user_created)

## Running

```bash
# Locally
go run cmd/notification-service/main.go

# Docker
docker-compose up notification-service
```

## Health Check

```
GET /health
```

Response:
```json
{
  "status": "ok",
  "service": "notification-service",
  "timestamp": "2024-01-01T12:00:00Z"
}
```
