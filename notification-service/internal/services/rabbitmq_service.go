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

	// Create connection with automatic reconnection
	conn, err := rabbitmq.NewConn(
		cfg.RabbitMQ.URL,
		rabbitmq.WithConnectionOptionsLogging,
		rabbitmq.WithConnectionOptionsReconnectInterval(5), // 5 seconds between reconnection attempts
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Create consumer with automatic reconnection
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
	err := r.consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		// Check context for graceful shutdown
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

	// Create notification in database
	notification := &models.Notification{
		UserID:  event.UserID,
		Type:    "user.created",
		Message: fmt.Sprintf("Welcome! Your account has been successfully created for email: %s", event.Email),
		Status:  "pending",
	}

	if err := r.db.Create(notification).Error; err != nil {
		return fmt.Errorf("failed to create notification record: %v", err)
	}

	// TODO: Add email/SMS sending logic here
	log.Printf("Would send welcome email to: %s", event.Email)

	return nil
}

func (r *RabbitMQService) Close() {
	// Cancel context for graceful shutdown
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
