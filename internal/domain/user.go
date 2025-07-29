package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Phone     string         `json:"phone" gorm:"index" validate:"omitempty,min=10,max=15"`
	Status    UserStatus     `json:"status" gorm:"default:active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

func (u UserStatus) IsValid() bool {
	switch u {
	case UserStatusActive, UserStatusInactive, UserStatusSuspended:
		return true
	}
	return false
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List(limit, offset int) ([]*User, error)
	Count() (int64, error)
}

type UserMessage struct {
	Action    string `json:"action"`
	UserID    uint   `json:"user_id"`
	User      *User  `json:"user,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	UserActionCreated = "user.created"
	UserActionUpdated = "user.updated"
	UserActionDeleted = "user.deleted"
)