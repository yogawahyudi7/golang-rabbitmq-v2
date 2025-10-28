package rabbitmq

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang-rabbitmq-v2/pkg/logger"
)

// EventType defines the type of event
type EventType string

const (
	EventPublish      EventType = "publish"
	EventConsume      EventType = "consume"
	EventConfirmAck   EventType = "confirm_ack"
	EventConfirmNack  EventType = "confirm_nack"
	EventConfirmTimeout EventType = "confirm_timeout"
	EventFailed       EventType = "failed"
	EventRequeued     EventType = "requeued"
)

// Event represents a single pub/sub event
type Event struct {
	ID          string                 `json:"id"`
	Type        EventType              `json:"type"`
	Timestamp   time.Time              `json:"timestamp"`
	RoutingKey  string                 `json:"routing_key,omitempty"`
	DeliveryTag uint64                 `json:"delivery_tag,omitempty"`
	Duration    time.Duration          `json:"duration_ms,omitempty"`
	Success     bool                   `json:"success"`
	Error       string                 `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ActivitySnapshot represents current ongoing activities
type ActivitySnapshot struct {
	Timestamp          time.Time              `json:"timestamp"`
	ActivePublishers   int                    `json:"active_publishers"`
	ActiveConsumers    int                    `json:"active_consumers"`
	PublishingInProgress int                  `json:"publishing_in_progress"`
	ConsumingInProgress  int                  `json:"consuming_in_progress"`
	CurrentPublishRate float64               `json:"current_publish_rate"`
	CurrentConsumeRate float64               `json:"current_consume_rate"`
	RecentEvents       []Event                `json:"recent_events"`
}

// Metrics menyimpan statistik internal untuk monitoring RabbitMQ
type Metrics struct {
	// Message counters
	messagesPublished   int64 // Total message yang berhasil dipublish
	messagesConsumed    int64 // Total message yang berhasil diproses
	messagesFailed      int64 // Total message yang gagal diproses
	messagesRequeued    int64 // Total message yang di-requeue

	// Publisher confirm counters
	confirmsReceived    int64 // Total confirmation yang diterima
	confirmsAcked       int64 // Total message yang di-ack broker
	confirmsNacked      int64 // Total message yang di-nack broker
	confirmsTimeout     int64 // Total confirmation yang timeout

	// Connection metrics
	connectionsCreated  int64 // Total koneksi yang dibuat
	reconnectAttempts   int64 // Total attempt reconnection
	connectionErrors    int64 // Total error koneksi

	// Processing metrics
	processingDuration  int64 // Total waktu processing (nanoseconds)
	averageProcessTime  int64 // Average processing time (nanoseconds)

	// Rate tracking
	publishRate         *RateTracker // Publish rate per detik
	consumeRate         *RateTracker // Consume rate per detik

	// Activity tracking
	activePublishers    int64 // Current active publishers
	activeConsumers     int64 // Current active consumers
	publishingInProgress int64 // Current publishing operations
	consumingInProgress  int64 // Current consuming operations

	// Event history
	events              []Event // History of events
	eventsMu            sync.RWMutex
	maxEvents           int     // Maximum events to keep in history
	eventCounter        int64   // Counter for generating event IDs

	logger              *logger.Logger
	mu                  sync.RWMutex
	startTime           time.Time
}

// RateTracker untuk menghitung rate per detik
type RateTracker struct {
	count     int64     // Total count
	lastCount int64     // Count pada interval sebelumnya
	rate      float64   // Current rate per second
	lastCheck time.Time // Waktu pengecekan terakhir
	mu        sync.RWMutex
}

// NewMetrics membuat instance metrics baru
func NewMetrics(log *logger.Logger) *Metrics {
	return &Metrics{
		publishRate: NewRateTracker(),
		consumeRate: NewRateTracker(),
		logger:      log,
		startTime:   time.Now(),
		events:      make([]Event, 0, 1000),
		maxEvents:   1000, // Keep last 1000 events
	}
}

// NewRateTracker membuat rate tracker baru
func NewRateTracker() *RateTracker {
	return &RateTracker{
		lastCheck: time.Now(),
	}
}

// === Publisher Metrics ===

// IncrementPublished menambah counter message yang berhasil dipublish
func (m *Metrics) IncrementPublished() {
	atomic.AddInt64(&m.messagesPublished, 1)
	m.publishRate.Increment()
}

// IncrementConfirmAcked menambah counter confirmation yang di-ack
func (m *Metrics) IncrementConfirmAcked() {
	atomic.AddInt64(&m.confirmsAcked, 1)
	atomic.AddInt64(&m.confirmsReceived, 1)
}

// IncrementConfirmNacked menambah counter confirmation yang di-nack
func (m *Metrics) IncrementConfirmNacked() {
	atomic.AddInt64(&m.confirmsNacked, 1)
	atomic.AddInt64(&m.confirmsReceived, 1)
}

// IncrementConfirmTimeout menambah counter confirmation timeout
func (m *Metrics) IncrementConfirmTimeout() {
	atomic.AddInt64(&m.confirmsTimeout, 1)
}

// === Consumer Metrics ===

// IncrementConsumed menambah counter message yang berhasil diproses
func (m *Metrics) IncrementConsumed() {
	atomic.AddInt64(&m.messagesConsumed, 1)
	m.consumeRate.Increment()
}

// IncrementFailed menambah counter message yang gagal diproses
func (m *Metrics) IncrementFailed() {
	atomic.AddInt64(&m.messagesFailed, 1)
}

// IncrementRequeued menambah counter message yang di-requeue
func (m *Metrics) IncrementRequeued() {
	atomic.AddInt64(&m.messagesRequeued, 1)
}

// RecordProcessingTime mencatat waktu processing message
func (m *Metrics) RecordProcessingTime(duration time.Duration) {
	nanos := duration.Nanoseconds()
	atomic.AddInt64(&m.processingDuration, nanos)
	
	// Update average processing time
	totalMessages := atomic.LoadInt64(&m.messagesConsumed)
	if totalMessages > 0 {
		totalDuration := atomic.LoadInt64(&m.processingDuration)
		avg := totalDuration / totalMessages
		atomic.StoreInt64(&m.averageProcessTime, avg)
	}
}

// === Connection Metrics ===

// IncrementConnectionCreated menambah counter koneksi yang dibuat
func (m *Metrics) IncrementConnectionCreated() {
	atomic.AddInt64(&m.connectionsCreated, 1)
}

// IncrementReconnectAttempt menambah counter reconnection attempt
func (m *Metrics) IncrementReconnectAttempt() {
	atomic.AddInt64(&m.reconnectAttempts, 1)
}

// IncrementConnectionError menambah counter error koneksi
func (m *Metrics) IncrementConnectionError() {
	atomic.AddInt64(&m.connectionErrors, 1)
}

// === Rate Tracker Methods ===

// Increment menambah counter untuk rate tracking
func (rt *RateTracker) Increment() {
	atomic.AddInt64(&rt.count, 1)
}

// GetRate menghitung dan return current rate per detik
func (rt *RateTracker) GetRate() float64 {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(rt.lastCheck).Seconds()
	
	// Update rate setiap detik
	if elapsed >= 1.0 {
		currentCount := atomic.LoadInt64(&rt.count)
		countDiff := currentCount - rt.lastCount
		
		rt.rate = float64(countDiff) / elapsed
		rt.lastCount = currentCount
		rt.lastCheck = now
	}
	
	return rt.rate
}

// === Getter Methods untuk Metrics ===

// GetPublishedCount return total message yang dipublish
func (m *Metrics) GetPublishedCount() int64 {
	return atomic.LoadInt64(&m.messagesPublished)
}

// GetConsumedCount return total message yang diproses
func (m *Metrics) GetConsumedCount() int64 {
	return atomic.LoadInt64(&m.messagesConsumed)
}

// GetFailedCount return total message yang gagal
func (m *Metrics) GetFailedCount() int64 {
	return atomic.LoadInt64(&m.messagesFailed)
}

// GetRequeuedCount return total message yang di-requeue
func (m *Metrics) GetRequeuedCount() int64 {
	return atomic.LoadInt64(&m.messagesRequeued)
}

// GetConfirmsAcked return total confirmation yang di-ack
func (m *Metrics) GetConfirmsAcked() int64 {
	return atomic.LoadInt64(&m.confirmsAcked)
}

// GetConfirmsNacked return total confirmation yang di-nack
func (m *Metrics) GetConfirmsNacked() int64 {
	return atomic.LoadInt64(&m.confirmsNacked)
}

// GetConfirmsTimeout return total confirmation timeout
func (m *Metrics) GetConfirmsTimeout() int64 {
	return atomic.LoadInt64(&m.confirmsTimeout)
}

// GetConnectionsCreated return total koneksi yang dibuat
func (m *Metrics) GetConnectionsCreated() int64 {
	return atomic.LoadInt64(&m.connectionsCreated)
}

// GetReconnectAttempts return total reconnection attempts
func (m *Metrics) GetReconnectAttempts() int64 {
	return atomic.LoadInt64(&m.reconnectAttempts)
}

// GetConnectionErrors return total connection errors
func (m *Metrics) GetConnectionErrors() int64 {
	return atomic.LoadInt64(&m.connectionErrors)
}

// GetAverageProcessingTime return rata-rata waktu processing
func (m *Metrics) GetAverageProcessingTime() time.Duration {
	nanos := atomic.LoadInt64(&m.averageProcessTime)
	return time.Duration(nanos)
}

// GetPublishRate return publish rate per detik
func (m *Metrics) GetPublishRate() float64 {
	return m.publishRate.GetRate()
}

// GetConsumeRate return consume rate per detik
func (m *Metrics) GetConsumeRate() float64 {
	return m.consumeRate.GetRate()
}

// GetUptime return waktu aplikasi berjalan
func (m *Metrics) GetUptime() time.Duration {
	return time.Since(m.startTime)
}

// === Summary Methods ===

// GetSummary return summary metrics dalam map
func (m *Metrics) GetSummary() map[string]interface{} {
	return map[string]interface{}{
		"messages": map[string]interface{}{
			"published": m.GetPublishedCount(),
			"consumed":  m.GetConsumedCount(),
			"failed":    m.GetFailedCount(),
			"requeued":  m.GetRequeuedCount(),
		},
		"confirmations": map[string]interface{}{
			"acked":   m.GetConfirmsAcked(),
			"nacked":  m.GetConfirmsNacked(),
			"timeout": m.GetConfirmsTimeout(),
		},
		"connections": map[string]interface{}{
			"created":           m.GetConnectionsCreated(),
			"reconnect_attempts": m.GetReconnectAttempts(),
			"errors":            m.GetConnectionErrors(),
		},
		"performance": map[string]interface{}{
			"avg_processing_time_ms": m.GetAverageProcessingTime().Milliseconds(),
			"publish_rate_per_sec":   m.GetPublishRate(),
			"consume_rate_per_sec":   m.GetConsumeRate(),
			"uptime_seconds":         m.GetUptime().Seconds(),
		},
	}
}

// LogMetrics log current metrics ke logger
func (m *Metrics) LogMetrics() {
	summary := m.GetSummary()
	m.logger.WithContext("rabbitmq-metrics").WithFields(summary).Info("Current metrics")
}

// StartPeriodicLogging memulai periodic logging metrics
func (m *Metrics) StartPeriodicLogging(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			m.LogMetrics()
		}
	}()

	m.logger.WithContext("rabbitmq-metrics").WithFields(map[string]interface{}{
		"interval_seconds": interval.Seconds(),
	}).Info("Started periodic metrics logging")
}

// === Event Tracking Methods ===

// generateEventID generates unique event ID
func (m *Metrics) generateEventID() string {
	id := atomic.AddInt64(&m.eventCounter, 1)
	return fmt.Sprintf("evt_%d_%d", time.Now().Unix(), id)
}

// RecordEvent records a new event to history
func (m *Metrics) RecordEvent(event Event) {
	m.eventsMu.Lock()
	defer m.eventsMu.Unlock()

	// Generate ID if not provided
	if event.ID == "" {
		event.ID = m.generateEventID()
	}

	// Add timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Add to events slice
	m.events = append(m.events, event)

	// Maintain max events limit (FIFO)
	if len(m.events) > m.maxEvents {
		m.events = m.events[1:]
	}
}

// GetEvents returns event history with optional filters
func (m *Metrics) GetEvents(eventType EventType, limit int) []Event {
	m.eventsMu.RLock()
	defer m.eventsMu.RUnlock()

	var filtered []Event

	// Filter by type if specified
	for i := len(m.events) - 1; i >= 0; i-- {
		event := m.events[i]
		if eventType == "" || event.Type == eventType {
			filtered = append(filtered, event)
			if limit > 0 && len(filtered) >= limit {
				break
			}
		}
	}

	return filtered
}

// GetRecentEvents returns most recent events
func (m *Metrics) GetRecentEvents(limit int) []Event {
	m.eventsMu.RLock()
	defer m.eventsMu.RUnlock()

	if limit <= 0 {
		limit = 100 // default limit
	}

	start := len(m.events) - limit
	if start < 0 {
		start = 0
	}

	// Return copy of recent events
	recent := make([]Event, len(m.events)-start)
	copy(recent, m.events[start:])

	// Reverse to get newest first
	for i := 0; i < len(recent)/2; i++ {
		j := len(recent) - 1 - i
		recent[i], recent[j] = recent[j], recent[i]
	}

	return recent
}

// ClearEvents clears event history
func (m *Metrics) ClearEvents() {
	m.eventsMu.Lock()
	defer m.eventsMu.Unlock()
	m.events = make([]Event, 0, m.maxEvents)
}

// === Activity Tracking Methods ===

// StartPublishing marks start of publishing operation
func (m *Metrics) StartPublishing() {
	atomic.AddInt64(&m.publishingInProgress, 1)
}

// EndPublishing marks end of publishing operation
func (m *Metrics) EndPublishing() {
	atomic.AddInt64(&m.publishingInProgress, -1)
}

// StartConsuming marks start of consuming operation
func (m *Metrics) StartConsuming() {
	atomic.AddInt64(&m.consumingInProgress, 1)
}

// EndConsuming marks end of consuming operation
func (m *Metrics) EndConsuming() {
	atomic.AddInt64(&m.consumingInProgress, -1)
}

// RegisterPublisher increments active publisher count
func (m *Metrics) RegisterPublisher() {
	atomic.AddInt64(&m.activePublishers, 1)
}

// UnregisterPublisher decrements active publisher count
func (m *Metrics) UnregisterPublisher() {
	atomic.AddInt64(&m.activePublishers, -1)
}

// RegisterConsumer increments active consumer count
func (m *Metrics) RegisterConsumer() {
	atomic.AddInt64(&m.activeConsumers, 1)
}

// UnregisterConsumer decrements active consumer count
func (m *Metrics) UnregisterConsumer() {
	atomic.AddInt64(&m.activeConsumers, -1)
}

// GetActivitySnapshot returns current activity snapshot
func (m *Metrics) GetActivitySnapshot(recentEventsLimit int) ActivitySnapshot {
	return ActivitySnapshot{
		Timestamp:            time.Now(),
		ActivePublishers:     int(atomic.LoadInt64(&m.activePublishers)),
		ActiveConsumers:      int(atomic.LoadInt64(&m.activeConsumers)),
		PublishingInProgress: int(atomic.LoadInt64(&m.publishingInProgress)),
		ConsumingInProgress:  int(atomic.LoadInt64(&m.consumingInProgress)),
		CurrentPublishRate:   m.GetPublishRate(),
		CurrentConsumeRate:   m.GetConsumeRate(),
		RecentEvents:         m.GetRecentEvents(recentEventsLimit),
	}
}

// === Enhanced Tracking Methods ===

// TrackPublishEvent tracks a publish event with details
func (m *Metrics) TrackPublishEvent(routingKey string, success bool, duration time.Duration, err error) {
	event := Event{
		Type:       EventPublish,
		RoutingKey: routingKey,
		Duration:   duration,
		Success:    success,
		Metadata: map[string]interface{}{
			"duration_ms": duration.Milliseconds(),
		},
	}

	if err != nil {
		event.Error = err.Error()
	}

	m.RecordEvent(event)
}

// TrackConsumeEvent tracks a consume event with details
func (m *Metrics) TrackConsumeEvent(deliveryTag uint64, routingKey string, success bool, duration time.Duration, err error) {
	event := Event{
		Type:        EventConsume,
		RoutingKey:  routingKey,
		DeliveryTag: deliveryTag,
		Duration:    duration,
		Success:     success,
		Metadata: map[string]interface{}{
			"duration_ms": duration.Milliseconds(),
		},
	}

	if err != nil {
		event.Error = err.Error()
	}

	m.RecordEvent(event)
}

// TrackConfirmEvent tracks a confirmation event
func (m *Metrics) TrackConfirmEvent(deliveryTag uint64, acked bool) {
	eventType := EventConfirmAck
	if !acked {
		eventType = EventConfirmNack
	}

	event := Event{
		Type:        eventType,
		DeliveryTag: deliveryTag,
		Success:     acked,
		Metadata: map[string]interface{}{
			"acked": acked,
		},
	}

	m.RecordEvent(event)
}