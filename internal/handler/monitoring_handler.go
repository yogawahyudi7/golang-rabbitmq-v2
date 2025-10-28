package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"golang-rabbitmq-v2/pkg/logger"
	"golang-rabbitmq-v2/pkg/rabbitmq"
)

type MonitoringHandler struct {
	rabbitConn *rabbitmq.Connection
	logger     *logger.Logger
}

func NewMonitoringHandler(rabbitConn *rabbitmq.Connection, log *logger.Logger) *MonitoringHandler {
	return &MonitoringHandler{
		rabbitConn: rabbitConn,
		logger:     log,
	}
}

// GetHealthCheck returns health status of the application
func (h *MonitoringHandler) GetHealthCheck(c *gin.Context) {
	metrics := h.rabbitConn.GetMetrics()

	health := map[string]interface{}{
		"status": "healthy",
		"timestamp": metrics.GetUptime().String(),
		"uptime_seconds": metrics.GetUptime().Seconds(),
		"rabbitmq": map[string]interface{}{
			"connected": h.rabbitConn.IsConnected(),
			"messages_published": metrics.GetPublishedCount(),
			"messages_consumed": metrics.GetConsumedCount(),
		},
	}

	c.JSON(http.StatusOK, health)
}

// GetMetrics returns current RabbitMQ metrics
func (h *MonitoringHandler) GetMetrics(c *gin.Context) {
	metrics := h.rabbitConn.GetMetrics()
	summary := metrics.GetSummary()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// GetActivitySnapshot returns real-time activity snapshot
// GET /api/v1/monitoring/activity?limit=50
func (h *MonitoringHandler) GetActivitySnapshot(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	metrics := h.rabbitConn.GetMetrics()
	snapshot := metrics.GetActivitySnapshot(limit)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    snapshot,
	})
}

// GetEventHistory returns event history with optional filters
// GET /api/v1/monitoring/events?type=publish&limit=100
func (h *MonitoringHandler) GetEventHistory(c *gin.Context) {
	eventType := c.Query("type")
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	metrics := h.rabbitConn.GetMetrics()
	var events []rabbitmq.Event

	if eventType != "" {
		events = metrics.GetEvents(rabbitmq.EventType(eventType), limit)
	} else {
		events = metrics.GetRecentEvents(limit)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"events": events,
			"count":  len(events),
			"filters": gin.H{
				"type":  eventType,
				"limit": limit,
			},
		},
	})
}

// GetPublishEvents returns publish event history
// GET /api/v1/monitoring/events/publish?limit=50
func (h *MonitoringHandler) GetPublishEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	metrics := h.rabbitConn.GetMetrics()
	events := metrics.GetEvents(rabbitmq.EventPublish, limit)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"events": events,
			"count":  len(events),
		},
	})
}

// GetConsumeEvents returns consume event history
// GET /api/v1/monitoring/events/consume?limit=50
func (h *MonitoringHandler) GetConsumeEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	metrics := h.rabbitConn.GetMetrics()
	events := metrics.GetEvents(rabbitmq.EventConsume, limit)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"events": events,
			"count":  len(events),
		},
	})
}

// GetFailedEvents returns failed event history
// GET /api/v1/monitoring/events/failed?limit=50
func (h *MonitoringHandler) GetFailedEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	metrics := h.rabbitConn.GetMetrics()
	events := metrics.GetEvents(rabbitmq.EventFailed, limit)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"events": events,
			"count":  len(events),
		},
	})
}

// ClearEventHistory clears all event history
// DELETE /api/v1/monitoring/events
func (h *MonitoringHandler) ClearEventHistory(c *gin.Context) {
	metrics := h.rabbitConn.GetMetrics()
	metrics.ClearEvents()

	h.logger.WithContext("monitoring").Info("Event history cleared")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event history cleared successfully",
	})
}

// GetDetailedMetrics returns detailed metrics with breakdown
// GET /api/v1/monitoring/metrics/detailed
func (h *MonitoringHandler) GetDetailedMetrics(c *gin.Context) {
	metrics := h.rabbitConn.GetMetrics()

	detailed := gin.H{
		"timestamp": metrics.GetUptime().String(),
		"uptime_seconds": metrics.GetUptime().Seconds(),
		"messages": gin.H{
			"published": metrics.GetPublishedCount(),
			"consumed":  metrics.GetConsumedCount(),
			"failed":    metrics.GetFailedCount(),
			"requeued":  metrics.GetRequeuedCount(),
			"success_rate": h.calculateSuccessRate(
				metrics.GetConsumedCount(),
				metrics.GetFailedCount(),
			),
		},
		"confirmations": gin.H{
			"acked":   metrics.GetConfirmsAcked(),
			"nacked":  metrics.GetConfirmsNacked(),
			"timeout": metrics.GetConfirmsTimeout(),
			"ack_rate": h.calculateAckRate(
				metrics.GetConfirmsAcked(),
				metrics.GetConfirmsNacked(),
			),
		},
		"connections": gin.H{
			"created":            metrics.GetConnectionsCreated(),
			"reconnect_attempts": metrics.GetReconnectAttempts(),
			"errors":             metrics.GetConnectionErrors(),
		},
		"performance": gin.H{
			"avg_processing_time_ms": metrics.GetAverageProcessingTime().Milliseconds(),
			"publish_rate_per_sec":   metrics.GetPublishRate(),
			"consume_rate_per_sec":   metrics.GetConsumeRate(),
		},
		"activity": gin.H{
			"active_publishers":      0, // Will be updated when integrated
			"active_consumers":       0, // Will be updated when integrated
			"publishing_in_progress": 0,
			"consuming_in_progress":  0,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    detailed,
	})
}

// Helper functions

func (h *MonitoringHandler) calculateSuccessRate(consumed, failed int64) float64 {
	total := consumed + failed
	if total == 0 {
		return 0
	}
	return float64(consumed) / float64(total) * 100
}

func (h *MonitoringHandler) calculateAckRate(acked, nacked int64) float64 {
	total := acked + nacked
	if total == 0 {
		return 0
	}
	return float64(acked) / float64(total) * 100
}
