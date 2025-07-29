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

type UserUseCase interface {
	CreateUser(user *domain.User) error
	GetUserByID(id uint) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id uint) error
	ListUsers(limit, offset int) ([]*domain.User, int64, error)
}

type userUseCase struct {
	userRepo  domain.UserRepository
	publisher rabbitmq.Publisher
	validator *validator.Validator
	logger    *logger.Logger
}

func NewUserUseCase(
	userRepo domain.UserRepository,
	publisher rabbitmq.Publisher,
	validator *validator.Validator,
	logger *logger.Logger,
) UserUseCase {
	return &userUseCase{
		userRepo:  userRepo,
		publisher: publisher,
		validator: validator,
		logger:    logger,
	}
}

func (uc *userUseCase) CreateUser(user *domain.User) error {
	if err := uc.validator.Validate(user); err != nil {
		uc.logger.WithContext("user-usecase").WithError(err).Error("User validation failed")
		return fmt.Errorf("validation failed: %w", err)
	}

	if !user.Status.IsValid() {
		user.Status = domain.UserStatusActive
	}

	if err := uc.userRepo.Create(user); err != nil {
		return err
	}

	if err := uc.publishUserEvent(domain.UserActionCreated, user); err != nil {
		uc.logger.WithContext("user-usecase").WithError(err).WithFields(map[string]interface{}{
			"user_id": user.ID,
			"action":  domain.UserActionCreated,
		}).Warn("Failed to publish user event")
	}

	uc.logger.WithContext("user-usecase").WithFields(map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User created successfully")

	return nil
}

func (uc *userUseCase) GetUserByID(id uint) (*domain.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) GetUserByEmail(email string) (*domain.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) UpdateUser(user *domain.User) error {
	if user.ID == 0 {
		return fmt.Errorf("invalid user ID")
	}

	if err := uc.validator.Validate(user); err != nil {
		uc.logger.WithContext("user-usecase").WithError(err).Error("User validation failed")
		return fmt.Errorf("validation failed: %w", err)
	}

	if !user.Status.IsValid() {
		return fmt.Errorf("invalid user status: %s", user.Status)
	}

	existingUser, err := uc.userRepo.GetByID(user.ID)
	if err != nil {
		return err
	}

	user.CreatedAt = existingUser.CreatedAt
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(user); err != nil {
		return err
	}

	if err := uc.publishUserEvent(domain.UserActionUpdated, user); err != nil {
		uc.logger.WithContext("user-usecase").WithError(err).WithFields(map[string]interface{}{
			"user_id": user.ID,
			"action":  domain.UserActionUpdated,
		}).Warn("Failed to publish user event")
	}

	uc.logger.WithContext("user-usecase").WithFields(map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User updated successfully")

	return nil
}

func (uc *userUseCase) DeleteUser(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid user ID")
	}

	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := uc.userRepo.Delete(id); err != nil {
		return err
	}

	if err := uc.publishUserEvent(domain.UserActionDeleted, user); err != nil {
		uc.logger.WithContext("user-usecase").WithError(err).WithFields(map[string]interface{}{
			"user_id": id,
			"action":  domain.UserActionDeleted,
		}).Warn("Failed to publish user event")
	}

	uc.logger.WithContext("user-usecase").WithFields(map[string]interface{}{
		"user_id": id,
		"email":   user.Email,
	}).Info("User deleted successfully")

	return nil
}

func (uc *userUseCase) ListUsers(limit, offset int) ([]*domain.User, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, err := uc.userRepo.List(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.userRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (uc *userUseCase) publishUserEvent(action string, user *domain.User) error {
	message := domain.UserMessage{
		Action:    action,
		UserID:    user.ID,
		User:      user,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal user message: %w", err)
	}

	routingKey := fmt.Sprintf("user.%s", action)
	return uc.publisher.PublishWithConfirm(context.Background(), routingKey, data)
}