package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Koshsky/subs-service/notification-service/internal/config"
	"github.com/Koshsky/subs-service/notification-service/internal/models"
	"github.com/wagslane/go-rabbitmq"
	"gorm.io/gorm"
)

type RabbitMQService struct {
	conn     *rabbitmq.Conn
	consumer *rabbitmq.Consumer
	config   *config.Config
	db       *gorm.DB
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewRabbitMQService(cfg *config.Config, db *gorm.DB) (*RabbitMQService, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Создаем соединение с автоматическим реконнектом
	conn, err := rabbitmq.NewConn(
		cfg.RabbitMQ.URL,
		rabbitmq.WithConnectionOptionsLogging,
		rabbitmq.WithConnectionOptionsReconnectInterval(5), // 5 секунд между попытками реконнекта
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Создаем consumer с автоматическим реконнектом
	consumer, err := rabbitmq.NewConsumer(
		conn,
		cfg.RabbitMQ.Queue,
		rabbitmq.WithConsumerOptionsRoutingKey("user.created"),
		rabbitmq.WithConsumerOptionsExchangeName(cfg.RabbitMQ.Exchange),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeKind("topic"),
		rabbitmq.WithConsumerOptionsExchangeDurable,
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsLogging,
	)
	if err != nil {
		conn.Close()
		cancel()
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	service := &RabbitMQService{
		conn:     conn,
		consumer: consumer,
		config:   cfg,
		db:       db,
		ctx:      ctx,
		cancel:   cancel,
	}

	return service, nil
}

func (r *RabbitMQService) StartConsuming() error {
	log.Printf("Started consuming messages from queue: %s", r.config.RabbitMQ.Queue)

	err := r.consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		// Проверяем контекст для graceful shutdown
		select {
		case <-r.ctx.Done():
			return rabbitmq.NackDiscard
		default:
		}

		if err := r.handleUserCreated(d.Body); err != nil {
			log.Printf("Error handling message: %v", err)
			return rabbitmq.NackRequeue
		}

		return rabbitmq.Ack
	})

	if err != nil {
		return fmt.Errorf("failed to start consumer: %v", err)
	}

	return nil
}

func (r *RabbitMQService) handleUserCreated(data []byte) error {
	var event models.UserCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal user created event: %v", err)
	}

	log.Printf("Received user created event: UserID=%s, Email=%s", event.UserID, event.Email)

	// Создаем уведомление в базе данных
	notification := &models.Notification{
		UserID:  event.UserID,
		Type:    "user.created",
		Message: fmt.Sprintf("Добро пожаловать! Ваш аккаунт был успешно создан для email: %s", event.Email),
		Status:  "pending",
	}

	if err := r.db.Create(notification).Error; err != nil {
		return fmt.Errorf("failed to create notification record: %v", err)
	}

	log.Printf("Created notification record with ID: %d for user: %s", notification.ID, event.UserID)

	// TODO: Здесь можно добавить логику отправки email/SMS
	log.Printf("Would send welcome email to: %s", event.Email)

	return nil
}

func (r *RabbitMQService) Close() {
	// Отменяем контекст для graceful shutdown
	if r.cancel != nil {
		r.cancel()
	}

	if r.consumer != nil {
		r.consumer.Close()
	}

	if r.conn != nil {
		r.conn.Close()
	}
}
