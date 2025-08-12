package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification represents a notification record
type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uuid.UUID      `json:"user_id" gorm:"not null;type:uuid"`
	Type      string         `json:"type" gorm:"not null"`
	Message   string         `json:"message" gorm:"not null"`
	Status    string         `json:"status" gorm:"default:'pending'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// UserCreatedEvent represents the user created event from RabbitMQ
type UserCreatedEvent struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}
