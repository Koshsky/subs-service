package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Koshsky/subs-service/auth-service/internal/config"
	"github.com/Koshsky/subs-service/auth-service/internal/models"
	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *config.Config
}

type UserCreatedEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func NewRabbitMQService(cfg *config.Config) (*RabbitMQService, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		cfg.RabbitMQ.Exchange, // name
		"topic",               // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	return &RabbitMQService{
		conn:    conn,
		channel: ch,
		config:  cfg,
	}, nil
}

func (r *RabbitMQService) PublishUserCreated(user *models.User) error {
	event := UserCreatedEvent{
		UserID: user.ID.String(),
		Email:  user.Email,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal user created event: %v", err)
	}

	err = r.channel.Publish(
		r.config.RabbitMQ.Exchange, // exchange
		"user.created",             // routing key
		false,                      // mandatory
		false,                      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish user created event: %v", err)
	}

	log.Printf("Published user.created event for user: %s", user.Email)
	return nil
}

func (r *RabbitMQService) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
