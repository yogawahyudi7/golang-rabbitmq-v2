package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"golang-rabbitmq-v2/internal/domain"
	"golang-rabbitmq-v2/internal/usecase"
	"golang-rabbitmq-v2/pkg/logger"
	"golang-rabbitmq-v2/pkg/rabbitmq"
)

type MessageHandler struct {
	messageUseCase usecase.MessageUseCase
	connection     *rabbitmq.Connection
	logger         *logger.Logger
}

type CreateMessageRequest struct {
	UserID     uint                `json:"user_id" validate:"required"`
	Content    string              `json:"content" validate:"required,min=1,max=5000"`
	Type       domain.MessageType  `json:"type" validate:"omitempty"`
	RoutingKey string              `json:"routing_key" validate:"omitempty"`
	Priority   int                 `json:"priority" validate:"omitempty,min=0,max=10"`
	MaxRetries int                 `json:"max_retries" validate:"omitempty,min=1,max=10"`
}

type MessageResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	Content     string              `json:"content"`
	Type        domain.MessageType  `json:"type"`
	Status      domain.MessageStatus `json:"status"`
	RoutingKey  string              `json:"routing_key"`
	Priority    int                 `json:"priority"`
	RetryCount  int                 `json:"retry_count"`
	MaxRetries  int                 `json:"max_retries"`
	ProcessedAt *string             `json:"processed_at"`
	FailedAt    *string             `json:"failed_at"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
	User        *UserResponse       `json:"user,omitempty"`
}

type MessageListResponse struct {
	Messages []MessageResponse `json:"messages"`
	Total    int64             `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

func NewMessageHandler(messageUseCase usecase.MessageUseCase, connection *rabbitmq.Connection, logger *logger.Logger) *MessageHandler {
	return &MessageHandler{
		messageUseCase: messageUseCase,
		connection:     connection,
		logger:         logger,
	}
}

func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithContext("message-handler").WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	message := &domain.Message{
		UserID:     req.UserID,
		Content:    req.Content,
		Type:       req.Type,
		RoutingKey: req.RoutingKey,
		Priority:   req.Priority,
		MaxRetries: req.MaxRetries,
	}

	if err := h.messageUseCase.CreateMessage(message); err != nil {
		h.logger.WithContext("message-handler").WithError(err).Error("Failed to create message")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create message",
			"details": err.Error(),
		})
		return
	}

	h.logger.WithContext("message-handler").WithFields(map[string]interface{}{
		"message_id": message.ID,
		"user_id":    message.UserID,
		"type":       message.Type,
	}).Info("Message created via API")

	c.JSON(http.StatusCreated, gin.H{
		"message": "Message created successfully",
		"data":    h.toMessageResponse(message),
	})
}

func (h *MessageHandler) GetMessage(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid message ID",
		})
		return
	}

	message, err := h.messageUseCase.GetMessageByID(uint(id))
	if err != nil {
		h.logger.WithContext("message-handler").WithError(err).WithFields(map[string]interface{}{
			"message_id": id,
		}).Error("Failed to get message")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Message not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": h.toMessageResponse(message),
	})
}

func (h *MessageHandler) ListMessages(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	offsetParam := c.DefaultQuery("offset", "0")
	userIDParam := c.Query("user_id")
	statusParam := c.Query("status")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	var messages []*domain.Message
	var total int64

	if userIDParam != "" {
		userID, err := strconv.ParseUint(userIDParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		messages, total, err = h.messageUseCase.ListMessagesByUser(uint(userID), limit, offset)
		if err != nil {
			h.logger.WithContext("message-handler").WithError(err).Error("Failed to list messages by user")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to list messages",
				"details": err.Error(),
			})
			return
		}
	} else if statusParam != "" {
		status := domain.MessageStatus(statusParam)
		if !status.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid message status",
			})
			return
		}

		messages, total, err = h.messageUseCase.ListMessagesByStatus(status, limit, offset)
		if err != nil {
			h.logger.WithContext("message-handler").WithError(err).Error("Failed to list messages by status")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to list messages",
				"details": err.Error(),
			})
			return
		}
	} else {
		messages, total, err = h.messageUseCase.ListMessages(limit, offset)
		if err != nil {
			h.logger.WithContext("message-handler").WithError(err).Error("Failed to list messages")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to list messages",
				"details": err.Error(),
			})
			return
		}
	}

	var messageResponses []MessageResponse
	for _, message := range messages {
		messageResponses = append(messageResponses, h.toMessageResponse(message))
	}

	response := MessageListResponse{
		Messages: messageResponses,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}

	c.JSON(http.StatusOK, response)
}

func (h *MessageHandler) RetryFailedMessages(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	if err := h.messageUseCase.RetryFailedMessages(limit); err != nil {
		h.logger.WithContext("message-handler").WithError(err).Error("Failed to retry failed messages")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retry failed messages",
			"details": err.Error(),
		})
		return
	}

	h.logger.WithContext("message-handler").WithFields(map[string]interface{}{
		"limit": limit,
	}).Info("Retry failed messages initiated via API")

	c.JSON(http.StatusOK, gin.H{
		"message": "Failed messages retry initiated successfully",
	})
}

func (h *MessageHandler) GetMessageStats(c *gin.Context) {
	stats, err := h.messageUseCase.GetMessageStats()
	if err != nil {
		h.logger.WithContext("message-handler").WithError(err).Error("Failed to get message stats")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get message stats",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

func (h *MessageHandler) GetRabbitMQMetrics(c *gin.Context) {
	metrics := h.connection.GetMetrics()
	summary := metrics.GetSummary()

	c.JSON(http.StatusOK, gin.H{
		"rabbitmq_metrics": summary,
	})
}

func (h *MessageHandler) GetHealthCheck(c *gin.Context) {
	isConnected := h.connection.IsConnected()
	
	status := "healthy"
	httpStatus := http.StatusOK
	
	if !isConnected {
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	metrics := h.connection.GetMetrics()
	summary := metrics.GetSummary()

	c.JSON(httpStatus, gin.H{
		"status":           status,
		"rabbitmq_connected": isConnected,
		"metrics":          summary,
		"timestamp":        metrics.GetUptime().String(),
	})
}

func (h *MessageHandler) toMessageResponse(message *domain.Message) MessageResponse {
	response := MessageResponse{
		ID:         message.ID,
		UserID:     message.UserID,
		Content:    message.Content,
		Type:       message.Type,
		Status:     message.Status,
		RoutingKey: message.RoutingKey,
		Priority:   message.Priority,
		RetryCount: message.RetryCount,
		MaxRetries: message.MaxRetries,
		CreatedAt:  message.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  message.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if message.ProcessedAt != nil {
		processedAt := message.ProcessedAt.Format("2006-01-02T15:04:05Z07:00")
		response.ProcessedAt = &processedAt
	}

	if message.FailedAt != nil {
		failedAt := message.FailedAt.Format("2006-01-02T15:04:05Z07:00")
		response.FailedAt = &failedAt
	}

	if message.User.ID != 0 {
		userResponse := UserResponse{
			ID:        message.User.ID,
			Name:      message.User.Name,
			Email:     message.User.Email,
			Phone:     message.User.Phone,
			Status:    message.User.Status,
			CreatedAt: message.User.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: message.User.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response.User = &userResponse
	}

	return response
}