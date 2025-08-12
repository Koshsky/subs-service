package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Koshsky/subs-service/auth-service/internal/config"
	"github.com/Koshsky/subs-service/auth-service/internal/models"
	"github.com/google/uuid"
	"github.com/wagslane/go-rabbitmq"
)

type RabbitMQService struct {
	conn      *rabbitmq.Conn
	publisher *rabbitmq.Publisher
	config    *config.Config
}

type UserCreatedEvent struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

func NewRabbitMQService(cfg *config.Config) (*RabbitMQService, error) {
	// Создаем соединение с автоматическим реконнектом
	conn, err := rabbitmq.NewConn(
		cfg.RabbitMQ.URL,
		rabbitmq.WithConnectionOptionsLogging,
		rabbitmq.WithConnectionOptionsReconnectInterval(5), // 5 секунд между попытками реконнекта
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Создаем publisher с автоматическим реконнектом
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(cfg.RabbitMQ.Exchange),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeKind("topic"),
		rabbitmq.WithPublisherOptionsExchangeDurable,
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create publisher: %v", err)
	}

	return &RabbitMQService{
		conn:      conn,
		publisher: publisher,
		config:    cfg,
	}, nil
}

func (r *RabbitMQService) PublishUserCreated(user *models.User) error {
	event := UserCreatedEvent{
		UserID: user.ID,
		Email:  user.Email,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal user created event: %v", err)
	}

	err = r.publisher.Publish(
		body,
		[]string{"user.created"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(r.config.RabbitMQ.Exchange),
	)
	if err != nil {
		return fmt.Errorf("failed to publish user created event: %v", err)
	}

	log.Printf("Published user.created event for user: %s", user.Email)
	return nil
}

func (r *RabbitMQService) Close() {
	if r.publisher != nil {
		r.publisher.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
