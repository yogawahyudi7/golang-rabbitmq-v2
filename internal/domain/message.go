package domain

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Content     string         `json:"content" gorm:"type:text;not null" validate:"required,min=1,max=5000"`
	Type        MessageType    `json:"type" gorm:"default:general"`
	Status      MessageStatus  `json:"status" gorm:"default:pending"`
	RoutingKey  string         `json:"routing_key" gorm:"index"`
	Priority    int            `json:"priority" gorm:"default:0"`
	RetryCount  int            `json:"retry_count" gorm:"default:0"`
	MaxRetries  int            `json:"max_retries" gorm:"default:3"`
	ProcessedAt *time.Time     `json:"processed_at"`
	FailedAt    *time.Time     `json:"failed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	User        User           `json:"user" gorm:"foreignKey:UserID"`
}

type MessageType string

const (
	MessageTypeGeneral      MessageType = "general"
	MessageTypeNotification MessageType = "notification"
	MessageTypeEmail        MessageType = "email"
	MessageTypeSMS          MessageType = "sms"
	MessageTypeAlert        MessageType = "alert"
)

func (m MessageType) IsValid() bool {
	switch m {
	case MessageTypeGeneral, MessageTypeNotification, MessageTypeEmail, MessageTypeSMS, MessageTypeAlert:
		return true
	}
	return false
}

type MessageStatus string

const (
	MessageStatusPending    MessageStatus = "pending"
	MessageStatusProcessing MessageStatus = "processing"
	MessageStatusProcessed  MessageStatus = "processed"
	MessageStatusFailed     MessageStatus = "failed"
	MessageStatusRetrying   MessageStatus = "retrying"
)

func (m MessageStatus) IsValid() bool {
	switch m {
	case MessageStatusPending, MessageStatusProcessing, MessageStatusProcessed, MessageStatusFailed, MessageStatusRetrying:
		return true
	}
	return false
}

type MessageRepository interface {
	Create(message *Message) error
	GetByID(id uint) (*Message, error)
	Update(message *Message) error
	Delete(id uint) error
	List(limit, offset int) ([]*Message, error)
	ListByUserID(userID uint, limit, offset int) ([]*Message, error)
	ListByStatus(status MessageStatus, limit, offset int) ([]*Message, error)
	Count() (int64, error)
	CountByStatus(status MessageStatus) (int64, error)
	GetFailedMessages(limit int) ([]*Message, error)
	GetRetryableMessages(limit int) ([]*Message, error)
}

type QueueMessage struct {
	MessageID uint        `json:"message_id"`
	UserID    uint        `json:"user_id"`
	Content   string      `json:"content"`
	Type      MessageType `json:"type"`
	Priority  int         `json:"priority"`
	Timestamp time.Time   `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}