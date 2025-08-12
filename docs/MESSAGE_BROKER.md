# Message Broker (RabbitMQ) - Документация

## Обзор

В данной системе используется **RabbitMQ** как брокер сообщений для реализации **event-driven архитектуры** между микросервисами. Это позволяет обеспечить асинхронное взаимодействие между сервисами и повысить надежность системы.

## Архитектура Message Broker

### Основные компоненты

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Auth Service  │───►│    RabbitMQ     │───►│Notification Svc │
│   (Publisher)   │    │ (Event Bus)     │    │  (Consumer)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Тип Exchange: Topic

В системе используется **Topic Exchange** - наиболее гибкий тип exchange в RabbitMQ, который позволяет:

- **Маршрутизировать сообщения** по routing key с использованием wildcards
- **Поддерживать сложные паттерны** маршрутизации
- **Обеспечивать масштабируемость** - легко добавлять новые consumers

### Конфигурация Exchange

```go
err = ch.ExchangeDeclare(
    cfg.RabbitMQ.Exchange, // name: "user_events"
    "topic",               // type: topic exchange
    true,                  // durable: переживает перезапуск RabbitMQ
    false,                 // auto-deleted: не удаляется автоматически
    false,                 // internal: доступен для публикации
    false,                 // no-wait: ждем подтверждения
    nil,                   // arguments: дополнительные параметры
)
```

## Реализация Publisher (Auth Service)

### Структура сервиса

```go
type RabbitMQService struct {
    conn    *amqp.Connection  // Соединение с RabbitMQ
    channel *amqp.Channel     // Канал для операций
    config  *config.Config    // Конфигурация
}
```

### Инициализация

```go
func NewRabbitMQService(cfg *config.Config) (*RabbitMQService, error) {
    // 1. Устанавливаем соединение
    conn, err := amqp.Dial(cfg.RabbitMQ.URL)
    
    // 2. Открываем канал
    ch, err := conn.Channel()
    
    // 3. Объявляем exchange
    err = ch.ExchangeDeclare(...)
    
    return &RabbitMQService{conn, ch, cfg}, nil
}
```

### Публикация событий

```go
func (r *RabbitMQService) PublishUserCreated(user *models.User) error {
    // 1. Создаем событие
    event := UserCreatedEvent{
        UserID: user.ID.String(),
        Email:  user.Email,
    }
    
    // 2. Сериализуем в JSON
    body, err := json.Marshal(event)
    
    // 3. Публикуем в exchange
    err = r.channel.Publish(
        r.config.RabbitMQ.Exchange, // exchange: "user_events"
        "user.created",             // routing key
        false,                      // mandatory: не требовать подтверждения
        false,                      // immediate: не требовать немедленной доставки
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}
```

### Особенности реализации Publisher

1. **Отсутствие очередей**: Publisher только публикует в exchange, не создает очереди
2. **Простая маршрутизация**: Используется фиксированный routing key `"user.created"`
3. **JSON сериализация**: События сериализуются в JSON для совместимости
4. **Логирование**: Каждая публикация логируется для отладки

## Реализация Consumer (Notification Service)

### Структура сервиса

```go
type RabbitMQService struct {
    conn     *amqp.Connection
    channel  *amqp.Channel
    config   *config.Config
    db       *gorm.DB                    // База данных для сохранения уведомлений
    handlers map[string]func([]byte) error // Обработчики событий
}
```

### Инициализация

```go
func NewRabbitMQService(cfg *config.Config, db *gorm.DB) (*RabbitMQService, error) {
    // 1. Устанавливаем соединение и канал
    conn, err := amqp.Dial(cfg.RabbitMQ.URL)
    ch, err := conn.Channel()
    
    // 2. Объявляем exchange (должен совпадать с publisher)
    err = ch.ExchangeDeclare(...)
    
    // 3. Объявляем очередь
    q, err := ch.QueueDeclare(
        cfg.RabbitMQ.Queue, // name: "user_created"
        true,               // durable: переживает перезапуск
        false,              // delete when unused: не удалять автоматически
        false,              // exclusive: не эксклюзивная
        false,              // no-wait: ждем подтверждения
        nil,                // arguments
    )
    
    // 4. Привязываем очередь к exchange
    err = ch.QueueBind(
        q.Name,                // queue name
        "user.created",        // routing key
        cfg.RabbitMQ.Exchange, // exchange
        false, nil,
    )
    
    // 5. Регистрируем обработчики
    service.handlers["user.created"] = func(data []byte) error {
        return service.handleUserCreated(data)
    }
}
```

### Потребление сообщений

```go
func (r *RabbitMQService) StartConsuming() error {
    // 1. Начинаем потребление
    msgs, err := r.channel.Consume(
        r.config.RabbitMQ.Queue, // queue
        "",                      // consumer: пустая строка = auto-generated
        false,                   // auto-ack: ручное подтверждение
        false,                   // exclusive: не эксклюзивный
        false,                   // no-local: получаем все сообщения
        false,                   // no-wait: ждем подтверждения
        nil,                     // args
    )
    
    // 2. Запускаем обработку в горутине
    go func() {
        for d := range msgs {
            routingKey := d.RoutingKey
            if handler, exists := r.handlers[routingKey]; exists {
                if err := handler(d.Body); err != nil {
                    log.Printf("Error handling message: %v", err)
                    d.Nack(false, true) // requeue при ошибке
                } else {
                    d.Ack(false) // подтверждаем успешную обработку
                }
            } else {
                log.Printf("No handler found for routing key: %s", routingKey)
                d.Nack(false, false) // не requeue неизвестных сообщений
            }
        }
    }()
}
```

### Обработка событий

```go
func (r *RabbitMQService) handleUserCreated(data []byte) error {
    // 1. Десериализуем событие
    var event models.UserCreatedEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return fmt.Errorf("failed to unmarshal user created event: %v", err)
    }
    
    // 2. Создаем уведомление в базе данных
    notification := &models.Notification{
        UserID:  event.UserID,
        Type:    "user.created",
        Message: fmt.Sprintf("Добро пожаловать! Ваш аккаунт был успешно создан для email: %s", event.Email),
        Status:  "pending",
    }
    
    if err := r.db.Create(notification).Error; err != nil {
        return fmt.Errorf("failed to create notification record: %v", err)
    }
    
    // 3. TODO: Отправка email/SMS
    log.Printf("Would send welcome email to: %s", event.Email)
    
    return nil
}
```

## Особенности реализации

### 1. Обработка ошибок

**Publisher (Auth Service):**
- Ошибки публикации логируются, но не прерывают основной поток
- При неудачной публикации пользователь все равно создается

**Consumer (Notification Service):**
- Используется ручное подтверждение (manual acknowledgment)
- При ошибке обработки сообщение возвращается в очередь (requeue)
- Неизвестные routing keys не requeue

### 2. Устойчивость к сбоям

- **Durable exchange и queue**: переживают перезапуск RabbitMQ
- **Connection management**: правильное закрытие соединений
- **Error handling**: детальное логирование ошибок

### 3. Масштабируемость

- **Topic exchange**: легко добавлять новые consumers
- **Routing keys**: гибкая маршрутизация сообщений
- **Handler pattern**: легко добавлять новые типы событий

### 4. Мониторинг

- **Подробное логирование**: каждый этап обработки логируется
- **Health checks**: проверка доступности RabbitMQ
- **Metrics**: можно добавить метрики для мониторинга

## Конфигурация

### Переменные окружения

```bash
# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_EXCHANGE=user_events
RABBITMQ_QUEUE=user_created
```

### Docker Compose

```yaml
rabbitmq:
  image: rabbitmq:3-management-alpine
  environment:
    RABBITMQ_DEFAULT_USER: ${RABBITMQ_USER:-guest}
    RABBITMQ_DEFAULT_PASS: ${RABBITMQ_PASSWORD:-guest}
  ports:
    - "5672:5672"      # AMQP
    - "15672:15672"    # Management UI
  volumes:
    - rabbitmq_data:/var/lib/rabbitmq
```

## Управление через Web UI

RabbitMQ предоставляет веб-интерфейс для мониторинга:
- **URL**: http://localhost:15672
- **Логин**: guest
- **Пароль**: guest

### Возможности:
- Просмотр очередей и их состояния
- Мониторинг сообщений
- Управление exchange и bindings
- Просмотр логов

## Лучшие практики

### 1. Обработка ошибок
- Всегда используйте ручное подтверждение в production
- Логируйте все ошибки для отладки
- Реализуйте retry механизм для критических операций

### 2. Производительность
- Используйте connection pooling для высоконагруженных систем
- Настройте prefetch count для оптимальной производительности
- Мониторьте размер очередей

### 3. Безопасность
- Используйте отдельные пользователей для разных сервисов
- Настройте SSL/TLS для production
- Ограничьте права доступа пользователей

### 4. Мониторинг
- Настройте алерты на размер очередей
- Мониторьте время обработки сообщений
- Отслеживайте количество dead letter сообщений

## Расширение системы

### Добавление новых событий

1. **В Publisher (Auth Service):**
```go
func (r *RabbitMQService) PublishUserDeleted(userID string) error {
    event := UserDeletedEvent{UserID: userID}
    body, _ := json.Marshal(event)
    
    return r.channel.Publish(
        r.config.RabbitMQ.Exchange,
        "user.deleted",
        false, false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}
```

2. **В Consumer (Notification Service):**
```go
service.handlers["user.deleted"] = func(data []byte) error {
    return service.handleUserDeleted(data)
}

func (r *RabbitMQService) handleUserDeleted(data []byte) error {
    // Обработка события удаления пользователя
}
```

### Добавление новых consumers

Можно легко добавить новые сервисы, которые будут обрабатывать те же события:

```go
// Новый сервис для аналитики
service.handlers["user.created"] = func(data []byte) error {
    return service.handleUserCreatedAnalytics(data)
}
```

## Заключение

Реализация брокера сообщений в данной системе обеспечивает:

- **Надежность**: устойчивость к сбоям и перезапускам
- **Масштабируемость**: легко добавлять новые consumers
- **Гибкость**: event-driven архитектура
- **Мониторинг**: подробное логирование и веб-интерфейс

Система готова к расширению и может легко адаптироваться к новым требованиям.
