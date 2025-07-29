package repository

import (
	"fmt"

	"golang-rabbitmq-v2/internal/database"
	"golang-rabbitmq-v2/internal/domain"
	"golang-rabbitmq-v2/pkg/logger"
)

type userRepository struct {
	db     *database.Database
	logger *logger.Logger
}

func NewUserRepository(db *database.Database, log *logger.Logger) domain.UserRepository {
	return &userRepository{
		db:     db,
		logger: log,
	}
}

func (r *userRepository) Create(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		r.logger.WithContext("user-repository").WithError(err).Error("Failed to create user")
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.WithContext("user-repository").WithFields(map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User created successfully")

	return nil
}

func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		r.logger.WithContext("user-repository").WithFields(map[string]interface{}{
			"user_id": id,
		}).WithError(err).Error("Failed to get user by ID")
		return nil, fmt.Errorf("failed to get user by ID %d: %w", id, err)
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		r.logger.WithContext("user-repository").WithFields(map[string]interface{}{
			"email": email,
		}).WithError(err).Error("Failed to get user by email")
		return nil, fmt.Errorf("failed to get user by email %s: %w", email, err)
	}

	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	if err := r.db.Save(user).Error; err != nil {
		r.logger.WithContext("user-repository").WithFields(map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
		}).WithError(err).Error("Failed to update user")
		return fmt.Errorf("failed to update user: %w", err)
	}

	r.logger.WithContext("user-repository").WithFields(map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User updated successfully")

	return nil
}

func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		r.logger.WithContext("user-repository").WithFields(map[string]interface{}{
			"user_id": id,
		}).WithError(err).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user with ID %d: %w", id, err)
	}

	r.logger.WithContext("user-repository").WithFields(map[string]interface{}{
		"user_id": id,
	}).Info("User deleted successfully")

	return nil
}

func (r *userRepository) List(limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		r.logger.WithContext("user-repository").WithFields(map[string]interface{}{
			"limit":  limit,
			"offset": offset,
		}).WithError(err).Error("Failed to list users")
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&domain.User{}).Count(&count).Error; err != nil {
		r.logger.WithContext("user-repository").WithError(err).Error("Failed to count users")
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}