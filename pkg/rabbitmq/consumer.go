package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang-rabbitmq-v2/pkg/logger"

	"github.com/streadway/amqp"
)

type MessageHandler func(ctx context.Context, body []byte) error

type Consumer interface {
	Start(ctx context.Context, handler MessageHandler) error
	Stop() error
}

type consumer struct {
	conn       *Connection
	logger     *logger.Logger
	mu         sync.RWMutex
	stopped    bool
	stopChan   chan struct{}
	workerPool chan struct{} // Semaphore untuk limit concurrent workers
}

// NewConsumer membuat consumer baru dengan worker pool untuk concurrent processing
// workerCount - jumlah maksimal worker yang bisa memproses message bersamaan
func NewConsumer(conn *Connection, log *logger.Logger, workerCount int) Consumer {
	return &consumer{
		conn:       conn,
		logger:     log,
		stopChan:   make(chan struct{}),
		workerPool: make(chan struct{}, workerCount), // Semaphore pattern
	}
}

// Start memulai consuming message dengan graceful shutdown support
// Menggunakan goroutine pool untuk memproses message secara concurrent
func (c *consumer) Start(ctx context.Context, handler MessageHandler) error {
	c.mu.Lock()
	if c.stopped {
		c.mu.Unlock()
		return fmt.Errorf("consumer already stopped")
	}
	c.mu.Unlock()

	// Dapatkan channel dari connection
	c.conn.mu.RLock()
	if c.conn.closed {
		c.conn.mu.RUnlock()
		return fmt.Errorf("connection is closed")
	}

	channel := c.conn.channel
	c.conn.mu.RUnlock()

	if channel == nil {
		if err := c.conn.reconnect(); err != nil {
			return err
		}
		c.conn.mu.RLock()
		channel = c.conn.channel
		c.conn.mu.RUnlock()
	}

	// channel.Consume - register consumer untuk menerima message
	// Parameters: queue, consumer-tag, auto-ack, exclusive, no-local, no-wait, arguments
	msgs, err := channel.Consume(
		c.conn.config.Queue, // nama queue
		"",                  // consumer tag (auto-generate)
		false,               // auto-ack = false (manual acknowledgment)
		false,               // exclusive = false (tidak eksklusif)
		false,               // no-local = false
		false,               // no-wait = false
		nil,                 // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	c.logger.WithContext("rabbitmq-consumer").Info("Consumer started")

	// Main consuming loop
	for {
		select {
		case <-ctx.Done():
			c.logger.WithContext("rabbitmq-consumer").Info("Consumer stopped due to context cancellation")
			return ctx.Err()

		case <-c.stopChan:
			c.logger.WithContext("rabbitmq-consumer").Info("Consumer stopped gracefully")
			return nil

		case msg, ok := <-msgs:
			if !ok {
				// Channel tertutup, coba reconnect
				c.logger.WithContext("rabbitmq-consumer").Warn("Message channel closed, attempting to reconnect")
				if err := c.conn.reconnect(); err != nil {
					return err
				}
				// Restart consuming setelah reconnect
				return c.Start(ctx, handler)
			}

			// Process message menggunakan worker pool (non-blocking)
			go c.processMessageWithWorkerPool(ctx, msg, handler)
		}
	}
}

// processMessageWithWorkerPool memproses message dengan worker pool pattern
// Menggunakan semaphore untuk membatasi jumlah concurrent workers
func (c *consumer) processMessageWithWorkerPool(ctx context.Context, msg amqp.Delivery, handler MessageHandler) {
	// Acquire worker dari pool (blocking jika pool penuh)
	select {
	case c.workerPool <- struct{}{}:
		// Worker acquired, lanjutkan processing
	case <-ctx.Done():
		// Context cancelled, reject message
		msg.Nack(false, true)
		return
	}

	// Release worker ke pool setelah processing selesai
	defer func() { <-c.workerPool }()

	// Process message dengan timeout dan recovery
	c.processMessage(ctx, msg, handler)
}

// processMessage memproses individual message dengan error handling dan recovery
// msg.Ack() - acknowledge message berhasil diproses
// msg.Nack() - negative acknowledge dengan requeue option
func (c *consumer) processMessage(ctx context.Context, msg amqp.Delivery, handler MessageHandler) {
	startTime := time.Now()

	// Track consuming start
	c.conn.metrics.StartConsuming()
	defer c.conn.metrics.EndConsuming()

	// Panic recovery
	defer func() {
		if r := recover(); r != nil {
			duration := time.Since(startTime)
			c.logger.WithContext("rabbitmq-consumer").WithFields(map[string]interface{}{
				"panic":        r,
				"delivery_tag": msg.DeliveryTag,
				"routing_key":  msg.RoutingKey,
			}).Error("Recovered from panic while processing message")

			// Track failed event
			c.conn.metrics.IncrementFailed()
			c.conn.metrics.TrackConsumeEvent(
				msg.DeliveryTag,
				msg.RoutingKey,
				false,
				duration,
				fmt.Errorf("panic: %v", r),
			)

			// Nack message dengan requeue untuk retry
			msg.Nack(false, true)
		}
	}()

	// Create context dengan timeout untuk processing
	processCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Execute handler
	if err := handler(processCtx, msg.Body); err != nil {
		duration := time.Since(startTime)
		c.logger.WithContext("rabbitmq-consumer").WithFields(map[string]interface{}{
			"error":        err.Error(),
			"delivery_tag": msg.DeliveryTag,
			"routing_key":  msg.RoutingKey,
			"duration_ms":  duration.Milliseconds(),
		}).Error("Handler returned error")

		// Track metrics
		c.conn.metrics.IncrementFailed()
		c.conn.metrics.TrackConsumeEvent(msg.DeliveryTag, msg.RoutingKey, false, duration, err)

		// Tentukan apakah message harus di-requeue atau di-reject
		if c.shouldRequeue(err) {
			// Nack dengan requeue untuk retry
			c.conn.metrics.IncrementRequeued()
			msg.Nack(false, true)
		} else {
			// Nack tanpa requeue (dead letter atau discard)
			msg.Nack(false, false)
		}
		return
	}

	// Handler berhasil, acknowledge message
	duration := time.Since(startTime)
	c.logger.WithContext("rabbitmq-consumer").WithFields(map[string]interface{}{
		"delivery_tag": msg.DeliveryTag,
		"routing_key":  msg.RoutingKey,
		"duration_ms":  duration.Milliseconds(),
	}).Debug("Message processed successfully")

	// Track metrics
	c.conn.metrics.IncrementConsumed()
	c.conn.metrics.RecordProcessingTime(duration)
	c.conn.metrics.TrackConsumeEvent(msg.DeliveryTag, msg.RoutingKey, true, duration, nil)

	// msg.Ack - acknowledge message berhasil diproses
	// Parameter: multiple (false = ack single message)
	if err := msg.Ack(false); err != nil {
		c.logger.WithContext("rabbitmq-consumer").WithError(err).Error("Failed to acknowledge message")
	}
}

// shouldRequeue menentukan apakah message harus di-requeue berdasarkan jenis error
// Implementasi logic untuk retry policy
func (c *consumer) shouldRequeue(err error) bool {
	// Contoh logic untuk menentukan requeue:
	// - Database connection error -> requeue
	// - Validation error -> tidak requeue
	// - Temporary network error -> requeue

	// Untuk contoh ini, kita requeue semua error
	// Dalam production, implementasikan logic yang lebih sophisticated
	return true
}

// Stop menghentikan consumer secara graceful
// Mengirim signal ke stopChan untuk menghentikan consuming loop
func (c *consumer) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return nil
	}

	c.stopped = true
	close(c.stopChan)

	c.logger.WithContext("rabbitmq-consumer").Info("Consumer stop signal sent")
	return nil
}
