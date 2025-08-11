package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Koshsky/subs-service/notification-service/internal/config"
	"github.com/Koshsky/subs-service/notification-service/internal/models"
	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	config   *config.Config
	handlers map[string]func([]byte) error
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

	// Declare queue
	q, err := ch.QueueDeclare(
		cfg.RabbitMQ.Queue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		q.Name,                // queue name
		"user.created",        // routing key
		cfg.RabbitMQ.Exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	service := &RabbitMQService{
		conn:     conn,
		channel:  ch,
		config:   cfg,
		handlers: make(map[string]func([]byte) error),
	}

	service.handlers["user.created"] = func(data []byte) error {
		return service.handleUserCreated(data)
	}

	return service, nil
}

func (r *RabbitMQService) StartConsuming() error {
	msgs, err := r.channel.Consume(
		r.config.RabbitMQ.Queue, // queue
		"",                      // consumer
		false,                   // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %v", err)
	}

	log.Printf("Started consuming messages from queue: %s", r.config.RabbitMQ.Queue)

	go func() {
		for d := range msgs {
			routingKey := d.RoutingKey
			if handler, exists := r.handlers[routingKey]; exists {
				if err := handler(d.Body); err != nil {
					log.Printf("Error handling message: %v", err)
					d.Nack(false, true) // requeue
				} else {
					d.Ack(false)
				}
			} else {
				log.Printf("No handler found for routing key: %s", routingKey)
				d.Nack(false, false) // don't requeue
			}
		}
	}()

	return nil
}

func (r *RabbitMQService) handleUserCreated(data []byte) error {
	var event models.UserCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal user created event: %v", err)
	}

	log.Printf("Received user created event: UserID=%s, Email=%s", event.UserID, event.Email)

	// TODO: Implement actual notification logic
	// For now, just log the event
	log.Printf("Would send welcome email to: %s", event.Email)

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
