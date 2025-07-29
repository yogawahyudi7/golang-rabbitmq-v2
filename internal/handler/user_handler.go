package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"golang-rabbitmq-v2/internal/domain"
	"golang-rabbitmq-v2/internal/usecase"
	"golang-rabbitmq-v2/pkg/logger"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
	logger      *logger.Logger
}

type CreateUserRequest struct {
	Name   string             `json:"name" validate:"required,min=2,max=100"`
	Email  string             `json:"email" validate:"required,email"`
	Phone  string             `json:"phone" validate:"omitempty,min=10,max=15"`
	Status domain.UserStatus  `json:"status" validate:"omitempty"`
}

type UpdateUserRequest struct {
	Name   string             `json:"name" validate:"required,min=2,max=100"`
	Email  string             `json:"email" validate:"required,email"`
	Phone  string             `json:"phone" validate:"omitempty,min=10,max=15"`
	Status domain.UserStatus  `json:"status" validate:"required"`
}

type UserResponse struct {
	ID     uint               `json:"id"`
	Name   string             `json:"name"`
	Email  string             `json:"email"`
	Phone  string             `json:"phone"`
	Status domain.UserStatus  `json:"status"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
}

type UserListResponse struct {
	Users []UserResponse     `json:"users"`
	Total int64              `json:"total"`
	Limit int                `json:"limit"`
	Offset int               `json:"offset"`
}

func NewUserHandler(userUseCase usecase.UserUseCase, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithContext("user-handler").WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	user := &domain.User{
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Status: req.Status,
	}

	if err := h.userUseCase.CreateUser(user); err != nil {
		h.logger.WithContext("user-handler").WithError(err).Error("Failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	h.logger.WithContext("user-handler").WithFields(map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User created via API")

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    h.toUserResponse(user),
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := h.userUseCase.GetUserByID(uint(id))
	if err != nil {
		h.logger.WithContext("user-handler").WithError(err).WithFields(map[string]interface{}{
			"user_id": id,
		}).Error("Failed to get user")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": h.toUserResponse(user),
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithContext("user-handler").WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	user := &domain.User{
		ID:     uint(id),
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Status: req.Status,
	}

	if err := h.userUseCase.UpdateUser(user); err != nil {
		h.logger.WithContext("user-handler").WithError(err).WithFields(map[string]interface{}{
			"user_id": id,
		}).Error("Failed to update user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"details": err.Error(),
		})
		return
	}

	h.logger.WithContext("user-handler").WithFields(map[string]interface{}{
		"user_id": id,
		"email":   user.Email,
	}).Info("User updated via API")

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    h.toUserResponse(user),
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	if err := h.userUseCase.DeleteUser(uint(id)); err != nil {
		h.logger.WithContext("user-handler").WithError(err).WithFields(map[string]interface{}{
			"user_id": id,
		}).Error("Failed to delete user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"details": err.Error(),
		})
		return
	}

	h.logger.WithContext("user-handler").WithFields(map[string]interface{}{
		"user_id": id,
	}).Info("User deleted via API")

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	offsetParam := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	users, total, err := h.userUseCase.ListUsers(limit, offset)
	if err != nil {
		h.logger.WithContext("user-handler").WithError(err).Error("Failed to list users")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list users",
			"details": err.Error(),
		})
		return
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, h.toUserResponse(user))
	}

	response := UserListResponse{
		Users:  userResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) toUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}