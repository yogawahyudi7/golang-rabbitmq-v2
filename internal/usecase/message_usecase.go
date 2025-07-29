package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang-rabbitmq-v2/internal/domain"
	"golang-rabbitmq-v2/pkg/logger"
	"golang-rabbitmq-v2/pkg/rabbitmq"
	"golang-rabbitmq-v2/pkg/validator"
)

type MessageUseCase interface {
	CreateMessage(message *domain.Message) error
	ProcessMessage(ctx context.Context, queueMessage *domain.QueueMessage) error
	GetMessageByID(id uint) (*domain.Message, error)
	UpdateMessageStatus(id uint, status domain.MessageStatus) error
	ListMessages(limit, offset int) ([]*domain.Message, int64, error)
	ListMessagesByUser(userID uint, limit, offset int) ([]*domain.Message, int64, error)
	ListMessagesByStatus(status domain.MessageStatus, limit, offset int) ([]*domain.Message, int64, error)
	RetryFailedMessages(limit int) error
	GetMessageStats() (map[string]int64, error)
}

type messageUseCase struct {
	messageRepo domain.MessageRepository
	userRepo    domain.UserRepository
	publisher   rabbitmq.Publisher
	validator   *validator.Validator
	logger      *logger.Logger
}

func NewMessageUseCase(
	messageRepo domain.MessageRepository,
	userRepo domain.UserRepository,
	publisher rabbitmq.Publisher,
	validator *validator.Validator,
	logger *logger.Logger,
) MessageUseCase {
	return &messageUseCase{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		publisher:   publisher,
		validator:   validator,
		logger:      logger,
	}
}

func (uc *messageUseCase) CreateMessage(message *domain.Message) error {
	if err := uc.validator.Validate(message); err != nil {
		uc.logger.WithContext("message-usecase").WithError(err).Error("Message validation failed")
		return fmt.Errorf("validation failed: %w", err)
	}

	if !message.Type.IsValid() {
		message.Type = domain.MessageTypeGeneral
	}

	if !message.Status.IsValid() {
		message.Status = domain.MessageStatusPending
	}

	if message.Priority < 0 {
		message.Priority = 0
	}
	if message.Priority > 10 {
		message.Priority = 10
	}

	if message.MaxRetries <= 0 {
		message.MaxRetries = 3
	}

	user, err := uc.userRepo.GetByID(message.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.Status != domain.UserStatusActive {
		return fmt.Errorf("user is not active")
	}

	if err := uc.messageRepo.Create(message); err != nil {
		return err
	}

	queueMessage := &domain.QueueMessage{
		MessageID: message.ID,
		UserID:    message.UserID,
		Content:   message.Content,
		Type:      message.Type,
		Priority:  message.Priority,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"routing_key": message.RoutingKey,
			"max_retries": message.MaxRetries,
		},
	}

	if err := uc.publishToQueue(queueMessage); err != nil {
		uc.logger.WithContext("message-usecase").WithError(err).WithFields(map[string]interface{}{
			"message_id": message.ID,
		}).Error("Failed to publish message to queue")
		
		message.Status = domain.MessageStatusFailed
		uc.messageRepo.Update(message)
		return fmt.Errorf("failed to publish message to queue: %w", err)
	}

	uc.logger.WithContext("message-usecase").WithFields(map[string]interface{}{
		"message_id": message.ID,
		"user_id":    message.UserID,
		"type":       message.Type,
	}).Info("Message created and queued successfully")

	return nil
}

func (uc *messageUseCase) ProcessMessage(ctx context.Context, queueMessage *domain.QueueMessage) error {
	message, err := uc.messageRepo.GetByID(queueMessage.MessageID)
	if err != nil {
		return fmt.Errorf("message not found: %w", err)
	}

	if message.Status == domain.MessageStatusProcessed {
		uc.logger.WithContext("message-usecase").WithFields(map[string]interface{}{
			"message_id": message.ID,
		}).Warn("Message already processed")
		return nil
	}

	message.Status = domain.MessageStatusProcessing
	if err := uc.messageRepo.Update(message); err != nil {
		return err
	}

	startTime := time.Now()
	err = uc.processMessageByType(ctx, message, queueMessage)
	processingTime := time.Since(startTime)

	if err != nil {
		message.Status = domain.MessageStatusFailed
		message.RetryCount++
		now := time.Now()
		message.FailedAt = &now

		uc.logger.WithContext("message-usecase").WithError(err).WithFields(map[string]interface{}{
			"message_id":      message.ID,
			"retry_count":     message.RetryCount,
			"processing_time": processingTime.Milliseconds(),
		}).Error("Message processing failed")

		if message.RetryCount < message.MaxRetries {
			message.Status = domain.MessageStatusRetrying
			uc.logger.WithContext("message-usecase").WithFields(map[string]interface{}{
				"message_id":  message.ID,
				"retry_count": message.RetryCount,
				"max_retries": message.MaxRetries,
			}).Info("Message marked for retry")
		}

		uc.messageRepo.Update(message)
		return err
	}

	message.Status = domain.MessageStatusProcessed
	now := time.Now()
	message.ProcessedAt = &now
	
	if err := uc.messageRepo.Update(message); err != nil {
		return err
	}

	uc.logger.WithContext("message-usecase").WithFields(map[string]interface{}{
		"message_id":      message.ID,
		"processing_time": processingTime.Milliseconds(),
	}).Info("Message processed successfully")

	return nil
}

func (uc *messageUseCase) processMessageByType(ctx context.Context, message *domain.Message, queueMessage *domain.QueueMessage) error {
	switch message.Type {
	case domain.MessageTypeNotification:
		return uc.processNotification(ctx, message, queueMessage)
	case domain.MessageTypeEmail:
		return uc.processEmail(ctx, message, queueMessage)
	case domain.MessageTypeSMS:
		return uc.processSMS(ctx, message, queueMessage)
	case domain.MessageTypeAlert:
		return uc.processAlert(ctx, message, queueMessage)
	default:
		return uc.processGeneral(ctx, message, queueMessage)
	}
}

func (uc *messageUseCase) processNotification(ctx context.Context, message *domain.Message, queueMessage *domain.QueueMessage) error {
	uc.logger.WithContext("message-usecase").WithFields(map[string]interface{}{
		"message_id": message.ID,
		"user_id":    message.UserID,
	}).Info("Processing notification message")
	
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (uc *messageUseCase) processEmail(ctx context.Context, message *domain.Message, queueMessage *domain.QueueMessage) error {
	uc.logger.WithContext("message-usecase").WithFields(map[string]interface{}{
		"message_id": message.ID,
		"user_id":    message.UserID,
	}).Info("Processing email message")
	
	time.Sleep(500 * time.Millisecond)
	return nil
}

func (uc *messageUseCase) processSMS(ctx context.Context, message *domain.Message, queueMessage *domain.QueueMessage) error {
	uc.logger.WithContext("message-usecase").WithFields(map[string]interface{}{
		"message_id": message.ID,
		"user_id":    message.UserID,
	}).Info("Processing SMS message")
	
	time.Sleep(200 * time.Millisecond)
	return nil
}

func (uc *messageUseCase) processAlert(ctx context.Context, message *domain.Message, queueMessage *domain.QueueMessage) error {
	uc.logger.WithContext("message-usecase").WithFields(map[string]interface{}{
		"message_id": message.ID,
		"user_id":    message.UserID,
	}).Info("Processing alert message")
	
	time.Sleep(50 * time.Millisecond)
	return nil
}

func (uc *messageUseCase) processGeneral(ctx context.Context, message *domain.Message, queueMessage *domain.QueueMessage) error {
	uc.logger.WithContext("message-usecase").WithFields(map[string]interface{}{
		"message_id": message.ID,
		"user_id":    message.UserID,
	}).Info("Processing general message")
	
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (uc *messageUseCase) GetMessageByID(id uint) (*domain.Message, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid message ID")
	}

	return uc.messageRepo.GetByID(id)
}

func (uc *messageUseCase) UpdateMessageStatus(id uint, status domain.MessageStatus) error {
	if id == 0 {
		return fmt.Errorf("invalid message ID")
	}

	if !status.IsValid() {
		return fmt.Errorf("invalid message status: %s", status)
	}

	message, err := uc.messageRepo.GetByID(id)
	if err != nil {
		return err
	}

	message.Status = status
	return uc.messageRepo.Update(message)
}

func (uc *messageUseCase) ListMessages(limit, offset int) ([]*domain.Message, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	messages, err := uc.messageRepo.List(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.messageRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (uc *messageUseCase) ListMessagesByUser(userID uint, limit, offset int) ([]*domain.Message, int64, error) {
	if userID == 0 {
		return nil, 0, fmt.Errorf("invalid user ID")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	messages, err := uc.messageRepo.ListByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.messageRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (uc *messageUseCase) ListMessagesByStatus(status domain.MessageStatus, limit, offset int) ([]*domain.Message, int64, error) {
	if !status.IsValid() {
		return nil, 0, fmt.Errorf("invalid message status: %s", status)
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	messages, err := uc.messageRepo.ListByStatus(status, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.messageRepo.CountByStatus(status)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (uc *messageUseCase) RetryFailedMessages(limit int) error {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	messages, err := uc.messageRepo.GetRetryableMessages(limit)
	if err != nil {
		return err
	}

	retriedCount := 0
	for _, message := range messages {
		queueMessage := &domain.QueueMessage{
			MessageID: message.ID,
			UserID:    message.UserID,
			Content:   message.Content,
			Type:      message.Type,
			Priority:  message.Priority,
			Timestamp: time.Now(),
			Metadata: map[string]interface{}{
				"routing_key": message.RoutingKey,
				"max_retries": message.MaxRetries,
				"retry":       true,
			},
		}

		if err := uc.publishToQueue(queueMessage); err != nil {
			uc.logger.WithContext("message-usecase").WithError(err).WithFields(map[string]interface{}{
				"message_id": message.ID,
			}).Error("Failed to retry message")
			continue
		}

		message.Status = domain.MessageStatusPending
		if err := uc.messageRepo.Update(message); err != nil {
			uc.logger.WithContext("message-usecase").WithError(err).WithFields(map[string]interface{}{
				"message_id": message.ID,
			}).Error("Failed to update message status for retry")
			continue
		}

		retriedCount++
	}

	uc.logger.WithContext("message-usecase").WithFields(map[string]interface{}{
		"total_messages": len(messages),
		"retried_count":  retriedCount,
	}).Info("Completed retry failed messages")

	return nil
}

func (uc *messageUseCase) GetMessageStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	total, err := uc.messageRepo.Count()
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	pending, err := uc.messageRepo.CountByStatus(domain.MessageStatusPending)
	if err != nil {
		return nil, err
	}
	stats["pending"] = pending

	processing, err := uc.messageRepo.CountByStatus(domain.MessageStatusProcessing)
	if err != nil {
		return nil, err
	}
	stats["processing"] = processing

	processed, err := uc.messageRepo.CountByStatus(domain.MessageStatusProcessed)
	if err != nil {
		return nil, err
	}
	stats["processed"] = processed

	failed, err := uc.messageRepo.CountByStatus(domain.MessageStatusFailed)
	if err != nil {
		return nil, err
	}
	stats["failed"] = failed

	retrying, err := uc.messageRepo.CountByStatus(domain.MessageStatusRetrying)
	if err != nil {
		return nil, err
	}
	stats["retrying"] = retrying

	return stats, nil
}

func (uc *messageUseCase) publishToQueue(queueMessage *domain.QueueMessage) error {
	data, err := json.Marshal(queueMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal queue message: %w", err)
	}

	routingKey := fmt.Sprintf("message.%s", queueMessage.Type)
	if queueMessage.Priority >= 5 {
		routingKey += ".high"
	} else {
		routingKey += ".normal"
	}

	return uc.publisher.PublishWithConfirm(context.Background(), routingKey, data)
}