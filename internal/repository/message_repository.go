package repository

import (
	"fmt"
	"time"

	"golang-rabbitmq-v2/internal/database"
	"golang-rabbitmq-v2/internal/domain"
	"golang-rabbitmq-v2/pkg/logger"
)

type messageRepository struct {
	db     *database.Database
	logger *logger.Logger
}

func NewMessageRepository(db *database.Database, log *logger.Logger) domain.MessageRepository {
	return &messageRepository{
		db:     db,
		logger: log,
	}
}

func (r *messageRepository) Create(message *domain.Message) error {
	if err := r.db.Create(message).Error; err != nil {
		r.logger.WithContext("message-repository").WithError(err).Error("Failed to create message")
		return fmt.Errorf("failed to create message: %w", err)
	}

	r.logger.WithContext("message-repository").WithFields(map[string]interface{}{
		"message_id": message.ID,
		"user_id":    message.UserID,
		"type":       message.Type,
	}).Info("Message created successfully")

	return nil
}

func (r *messageRepository) GetByID(id uint) (*domain.Message, error) {
	var message domain.Message
	if err := r.db.Preload("User").First(&message, id).Error; err != nil {
		r.logger.WithContext("message-repository").WithFields(map[string]interface{}{
			"message_id": id,
		}).WithError(err).Error("Failed to get message by ID")
		return nil, fmt.Errorf("failed to get message by ID %d: %w", id, err)
	}

	return &message, nil
}

func (r *messageRepository) Update(message *domain.Message) error {
	if err := r.db.Save(message).Error; err != nil {
		r.logger.WithContext("message-repository").WithFields(map[string]interface{}{
			"message_id": message.ID,
			"status":     message.Status,
		}).WithError(err).Error("Failed to update message")
		return fmt.Errorf("failed to update message: %w", err)
	}

	r.logger.WithContext("message-repository").WithFields(map[string]interface{}{
		"message_id": message.ID,
		"status":     message.Status,
	}).Info("Message updated successfully")

	return nil
}

func (r *messageRepository) Delete(id uint) error {
	if err := r.db.Delete(&domain.Message{}, id).Error; err != nil {
		r.logger.WithContext("message-repository").WithFields(map[string]interface{}{
			"message_id": id,
		}).WithError(err).Error("Failed to delete message")
		return fmt.Errorf("failed to delete message with ID %d: %w", id, err)
	}

	r.logger.WithContext("message-repository").WithFields(map[string]interface{}{
		"message_id": id,
	}).Info("Message deleted successfully")

	return nil
}

func (r *messageRepository) List(limit, offset int) ([]*domain.Message, error) {
	var messages []*domain.Message
	if err := r.db.Preload("User").Limit(limit).Offset(offset).Order("created_at DESC").Find(&messages).Error; err != nil {
		r.logger.WithContext("message-repository").WithFields(map[string]interface{}{
			"limit":  limit,
			"offset": offset,
		}).WithError(err).Error("Failed to list messages")
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	return messages, nil
}

func (r *messageRepository) ListByUserID(userID uint, limit, offset int) ([]*domain.Message, error) {
	var messages []*domain.Message
	if err := r.db.Preload("User").Where("user_id = ?", userID).Limit(limit).Offset(offset).Order("created_at DESC").Find(&messages).Error; err != nil {
		r.logger.WithContext("message-repository").WithFields(map[string]interface{}{
			"user_id": userID,
			"limit":   limit,
			"offset":  offset,
		}).WithError(err).Error("Failed to list messages by user ID")
		return nil, fmt.Errorf("failed to list messages by user ID %d: %w", userID, err)
	}

	return messages, nil
}

func (r *messageRepository) ListByStatus(status domain.MessageStatus, limit, offset int) ([]*domain.Message, error) {
	var messages []*domain.Message
	if err := r.db.Preload("User").Where("status = ?", status).Limit(limit).Offset(offset).Order("created_at DESC").Find(&messages).Error; err != nil {
		r.logger.WithContext("message-repository").WithFields(map[string]interface{}{
			"status": status,
			"limit":  limit,
			"offset": offset,
		}).WithError(err).Error("Failed to list messages by status")
		return nil, fmt.Errorf("failed to list messages by status %s: %w", status, err)
	}

	return messages, nil
}

func (r *messageRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&domain.Message{}).Count(&count).Error; err != nil {
		r.logger.WithContext("message-repository").WithError(err).Error("Failed to count messages")
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return count, nil
}

func (r *messageRepository) CountByStatus(status domain.MessageStatus) (int64, error) {
	var count int64
	if err := r.db.Model(&domain.Message{}).Where("status = ?", status).Count(&count).Error; err != nil {
		r.logger.WithContext("message-repository").WithFields(map[string]interface{}{
			"status": status,
		}).WithError(err).Error("Failed to count messages by status")
		return 0, fmt.Errorf("failed to count messages by status %s: %w", status, err)
	}

	return count, nil
}

func (r *messageRepository) GetFailedMessages(limit int) ([]*domain.Message, error) {
	var messages []*domain.Message
	if err := r.db.Preload("User").Where("status = ? AND retry_count < max_retries", domain.MessageStatusFailed).Limit(limit).Order("created_at ASC").Find(&messages).Error; err != nil {
		r.logger.WithContext("message-repository").WithFields(map[string]interface{}{
			"limit": limit,
		}).WithError(err).Error("Failed to get failed messages")
		return nil, fmt.Errorf("failed to get failed messages: %w", err)
	}

	return messages, nil
}

func (r *messageRepository) GetRetryableMessages(limit int) ([]*domain.Message, error) {
	var messages []*domain.Message
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	
	if err := r.db.Preload("User").Where("status IN ? AND retry_count < max_retries AND (failed_at IS NULL OR failed_at < ?)", 
		[]domain.MessageStatus{domain.MessageStatusFailed, domain.MessageStatusRetrying}, fiveMinutesAgo).
		Limit(limit).Order("created_at ASC").Find(&messages).Error; err != nil {
		r.logger.WithContext("message-repository").WithFields(map[string]interface{}{
			"limit": limit,
		}).WithError(err).Error("Failed to get retryable messages")
		return nil, fmt.Errorf("failed to get retryable messages: %w", err)
	}

	return messages, nil
}