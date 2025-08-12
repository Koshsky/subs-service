# Поведение системы при проблемах с RabbitMQ

## Обзор

Данная система реализует **надежную event-driven архитектуру** с использованием библиотеки `go-rabbitmq`, которая обеспечивает **автоматическое переподключение** и **локальную очередь сообщений** для обеспечения максимальной надежности.

## Архитектура надежности

### Ключевые компоненты

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Auth Service  │───►│    RabbitMQ     │───►│Notification Svc │
│   (Publisher)   │    │ (Event Bus)     │    │  (Consumer)     │
│   + Local Queue │    │                 │    │ + Auto Reconnect│
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Технологии

- **go-rabbitmq**: Библиотека с встроенной поддержкой переподключения
- **Локальная очередь**: Буферизация сообщений при недоступности RabbitMQ
- **Exponential backoff**: Умная стратегия переподключения
- **Publisher confirms**: Гарантированная доставка сообщений

## Сценарии поведения

### Сценарий 1: Потеря соединения с RabbitMQ

#### Auth Service (Publisher)
```
1. RabbitMQ становится недоступен
2. Auth Service продолжает работать ✅
3. При создании пользователя:
   - Пользователь создается в БД ✅
   - Событие помещается в локальную очередь ✅
   - Логируется: "Message queued locally for retry"
4. Publisher worker пытается отправить сообщения каждые 5 секунд ✅
5. Сервис остается полностью функциональным ✅
```

#### Notification Service (Consumer)
```
1. RabbitMQ становится недоступен
2. Notification Service продолжает работать ✅
3. Consumer автоматически пытается переподключиться ✅
4. Используется exponential backoff для переподключения ✅
5. Health check показывает: "rabbitmq": "reconnecting", "consumer": "retrying"
6. Сервис остается стабильным ✅
```

### Сценарий 2: Восстановление RabbitMQ

#### Auth Service (Publisher)
```
1. RabbitMQ восстанавливается
2. Auth Service автоматически переподключается ✅
3. Все накопленные сообщения отправляются ✅
4. Логируется: "Successfully reconnected to RabbitMQ"
5. Логируется: "Flushed X messages from local queue"
6. Полная функциональность восстанавливается автоматически ✅
```

#### Notification Service (Consumer)
```
1. RabbitMQ восстанавливается
2. Notification Service автоматически переподключается ✅
3. Consumer возобновляет работу ✅
4. Health check показывает: "rabbitmq": "connected", "consumer": "running"
5. Полная функциональность восстанавливается автоматически ✅
```

## Конфигурация надежности

### Переменные окружения

```bash
# Настройки переподключения
RABBITMQ_RECONNECT_DELAY=5s                    # Начальная задержка
RABBITMQ_MAX_RECONNECT_ATTEMPTS=10             # Максимум попыток
RABBITMQ_PUBLISHER_BUFFER_SIZE=1000            # Размер локальной очереди
RABBITMQ_PUBLISHER_FLUSH_INTERVAL=5s           # Интервал отправки
```

### Стратегия переподключения

```go
// Exponential backoff с разумными лимитами
reconnectBackoff := backoff.NewExponentialBackOff()
reconnectBackoff.InitialInterval = 5 * time.Second
reconnectBackoff.MaxInterval = 30 * time.Second
reconnectBackoff.MaxElapsedTime = 0 // Без ограничения времени
```

## Логи и мониторинг

### Health Check Response

**При нормальной работе:**
```json
{
  "status": "ok",
  "service": "notification-service",
  "timestamp": "2024-01-15T10:30:00Z",
  "rabbitmq": "connected",
  "consumer": "running",
  "queue_size": 0
}
```

**При переподключении:**
```json
{
  "status": "ok",
  "service": "notification-service",
  "timestamp": "2024-01-15T10:30:00Z",
  "rabbitmq": "reconnecting",
  "consumer": "retrying",
  "reconnect_attempts": 3
}
```

**При использовании локальной очереди:**
```json
{
  "status": "ok",
  "service": "auth-service",
  "timestamp": "2024-01-15T10:30:00Z",
  "rabbitmq": "disconnected",
  "publisher": "buffering",
  "local_queue_size": 15
}
```

### Типичные логи

**Auth Service при потере RabbitMQ:**
```
2024/01/15 10:30:00 Warning: RabbitMQ connection lost, buffering messages locally
2024/01/15 10:30:01 User created successfully in database
2024/01/15 10:30:01 Message queued locally for retry (queue size: 1)
2024/01/15 10:30:05 Attempting to reconnect to RabbitMQ (attempt 1/10)
2024/01/15 10:30:05 Successfully reconnected to RabbitMQ
2024/01/15 10:30:05 Flushed 1 messages from local queue
```

**Notification Service при потере RabbitMQ:**
```
2024/01/15 10:30:00 Message channel closed, attempting to reconnect...
2024/01/15 10:30:05 Attempting to reconnect to RabbitMQ (attempt 1/10)
2024/01/15 10:30:05 Reconnection attempt 1 failed: connection refused
2024/01/15 10:30:10 Attempting to reconnect to RabbitMQ (attempt 2/10)
2024/01/15 10:30:10 Successfully reconnected to RabbitMQ after 2 attempts
2024/01/15 10:30:10 Consumer successfully reconnected
```

## Преимущества новой архитектуры

### ✅ Автоматическое восстановление
- **Нет необходимости в ручном вмешательстве**
- **Автоматическое переподключение с exponential backoff**
- **Восстановление всех функций без перезапуска**

### ✅ Сохранение данных
- **Локальная очередь предотвращает потерю сообщений**
- **Автоматическая отправка накопленных сообщений**
- **Retry логика с ограничением попыток**

### ✅ Надежность
- **Publisher confirms для гарантированной доставки**
- **Persistent delivery для сохранения сообщений**
- **Graceful shutdown с очисткой ресурсов**

### ✅ Масштабируемость
- **Многопоточная обработка сообщений**
- **Конфигурируемые параметры производительности**
- **Легкое добавление новых типов событий**

## Рекомендации по эксплуатации

### 1. Мониторинг

Настройте алерты на:
- `rabbitmq: "disconnected"` в health checks
- `local_queue_size > 100` (высокое количество накопленных сообщений)
- `reconnect_attempts > 5` (много попыток переподключения)

### 2. Настройка параметров

**Для высоконагруженных систем:**
```bash
RABBITMQ_PUBLISHER_BUFFER_SIZE=5000
RABBITMQ_PUBLISHER_FLUSH_INTERVAL=2s
RABBITMQ_MAX_RECONNECT_ATTEMPTS=20
```

**Для систем с нестабильным соединением:**
```bash
RABBITMQ_RECONNECT_DELAY=10s
RABBITMQ_MAX_RECONNECT_ATTEMPTS=30
```

### 3. Стратегии для production

#### Мониторинг и алерты
```yaml
# Prometheus metrics
rabbitmq_connection_status: 1/0
rabbitmq_local_queue_size: gauge
rabbitmq_reconnect_attempts: counter
rabbitmq_messages_published: counter
rabbitmq_messages_consumed: counter
```

#### Автоматическое восстановление
- **Kubernetes**: Используйте liveness/readiness probes
- **Docker Swarm**: Настройте restart policies
- **Systemd**: Используйте Restart=always

## Тестирование

### Запуск тестов отказоустойчивости

```bash
# Полный тест всех сценариев
./scripts/test-rabbitmq-resilience.sh

# Ручное тестирование
docker-compose stop rabbitmq
# Проверить health checks
curl http://localhost:8082/health
# Создать пользователя (должно работать с локальной очередью)
curl -X POST http://localhost:8081/api/auth/register -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"testpass"}'
# Запустить RabbitMQ обратно
docker-compose start rabbitmq
# Проверить, что сообщения отправились автоматически
curl http://localhost:8082/health
```

### Сценарии тестирования

1. **Кратковременная недоступность RabbitMQ** (30 секунд)
2. **Длительная недоступность** (5 минут)
3. **Множественные переподключения**
4. **Переполнение локальной очереди**
5. **Graceful shutdown во время сбоя**

## Заключение

Новая архитектура обеспечивает:

- **Максимальную надежность**: автоматическое восстановление без потери данных
- **Нулевое время простоя**: сервисы продолжают работать при сбоях RabbitMQ
- **Простоту эксплуатации**: нет необходимости в ручном вмешательстве
- **Масштабируемость**: готовность к высоким нагрузкам

Система теперь полностью готова для production использования с минимальными требованиями к мониторингу и обслуживанию.
