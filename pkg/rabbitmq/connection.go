package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"

	"golang-rabbitmq-v2/internal/config"
	"golang-rabbitmq-v2/pkg/logger"
)

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *config.RabbitMQConfig
	logger  *logger.Logger
	metrics *Metrics
	mu      sync.RWMutex
	closed  bool
}

// NewConnection membuat koneksi baru ke RabbitMQ
// Menginisialisasi struct Connection dan memanggil connect()
func NewConnection(cfg *config.RabbitMQConfig, log *logger.Logger) (*Connection, error) {
	conn := &Connection{
		config:  cfg,
		logger:  log,
		metrics: NewMetrics(log),
	}

	if err := conn.connect(); err != nil {
		return nil, err
	}

	return conn, nil
}

// connect membuat koneksi TCP ke RabbitMQ broker
// amqp.Dial() - membuat koneksi ke broker menggunakan AMQP URL
// conn.Channel() - membuat channel baru dari koneksi yang ada
func (c *Connection) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var err error
	// amqp.Dial - establish TCP connection ke RabbitMQ server
	c.conn, err = amqp.Dial(c.config.URL)
	if err != nil {
		c.metrics.IncrementConnectionError()
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	c.metrics.IncrementConnectionCreated()

	// conn.Channel - membuat channel baru untuk komunikasi
	c.channel, err = c.conn.Channel()
	if err != nil {
		c.conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	if err := c.setupExchangeAndQueue(); err != nil {
		c.channel.Close()
		c.conn.Close()
		return fmt.Errorf("failed to setup exchange and queue: %w", err)
	}

	c.logger.WithContext("rabbitmq").Info("Successfully connected to RabbitMQ")
	return nil
}

// setupExchangeAndQueue mendeklarasikan exchange, queue, dan binding
// ExchangeDeclare - membuat/memastikan exchange ada
// QueueDeclare - membuat/memastikan queue ada  
// QueueBind - mengikat queue ke exchange dengan routing key
func (c *Connection) setupExchangeAndQueue() error {
	// ExchangeDeclare - declare exchange dengan tipe "topic"
	// Parameters: name, type, durable, auto-delete, internal, no-wait, arguments
	err := c.channel.ExchangeDeclare(
		c.config.Exchange, // nama exchange
		"topic",           // tipe exchange (topic untuk routing key pattern)
		true,              // durable - exchange survive broker restart
		false,             // auto-delete - tidak delete otomatis
		false,             // internal - bisa diakses publisher
		false,             // no-wait - tunggu konfirmasi server
		nil,               // arguments tambahan
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// QueueDeclare - declare queue untuk menyimpan message
	// Parameters: name, durable, exclusive, auto-delete, no-wait, arguments
	_, err = c.channel.QueueDeclare(
		c.config.Queue, // nama queue
		true,           // durable - queue survive broker restart
		false,          // exclusive - bisa diakses multiple connection
		false,          // auto-delete - tidak delete otomatis
		false,          // no-wait - tunggu konfirmasi server
		nil,            // arguments tambahan
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// QueueBind - bind queue ke exchange dengan routing key
	// Parameters: queue, routing-key, exchange, no-wait, arguments
	err = c.channel.QueueBind(
		c.config.Queue,      // nama queue yang akan di-bind
		c.config.RoutingKey, // routing key pattern
		c.config.Exchange,   // nama exchange
		false,               // no-wait - tunggu konfirmasi server
		nil,                 // arguments tambahan
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

// reconnect melakukan reconnection otomatis dengan retry mechanism
// Dipanggil ketika koneksi terputus atau channel error
func (c *Connection) reconnect() error {
	c.logger.WithContext("rabbitmq").Warn("Attempting to reconnect to RabbitMQ")
	c.metrics.IncrementReconnectAttempt()
	
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if err := c.connect(); err != nil {
			c.logger.WithContext("rabbitmq").WithError(err).Errorf("Reconnection attempt %d failed", i+1)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		c.logger.WithContext("rabbitmq").Info("Successfully reconnected to RabbitMQ")
		return nil
	}
	
	return fmt.Errorf("failed to reconnect after %d attempts", maxRetries)
}

// Publish mengirim message ke exchange dengan routing key
// channel.Publish() - publish message ke exchange
func (c *Connection) Publish(ctx context.Context, routingKey string, body []byte) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return fmt.Errorf("connection is closed")
	}
	
	channel := c.channel
	c.mu.RUnlock()

	if channel == nil {
		if err := c.reconnect(); err != nil {
			return err
		}
		c.mu.RLock()
		channel = c.channel
		c.mu.RUnlock()
	}

	// channel.Publish - publish message ke exchange
	// Parameters: exchange, routing-key, mandatory, immediate, publishing
	err := channel.Publish(
		c.config.Exchange, // nama exchange tujuan
		routingKey,        // routing key untuk routing message
		false,             // mandatory - tidak return jika tidak ada queue
		false,             // immediate - tidak return jika tidak ada consumer
		amqp.Publishing{   // message properties
			ContentType:  "application/json", // tipe content
			Body:         body,               // isi message
			Timestamp:    time.Now(),         // timestamp pengiriman
			DeliveryMode: amqp.Persistent,    // persistent - survive broker restart
		},
	)

	if err != nil {
		c.logger.WithContext("rabbitmq").WithError(err).Error("Failed to publish message")
		return fmt.Errorf("failed to publish message: %w", err)
	}

	c.metrics.IncrementPublished()
	c.logger.WithContext("rabbitmq").WithFields(map[string]interface{}{
		"routing_key": routingKey,
		"body_size":   len(body),
	}).Debug("Message published successfully")

	return nil
}

// Consume memulai consuming message dari queue
// channel.Consume() - register consumer untuk menerima message
func (c *Connection) Consume(ctx context.Context, handler func([]byte) error) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return fmt.Errorf("connection is closed")
	}
	
	channel := c.channel
	c.mu.RUnlock()

	if channel == nil {
		if err := c.reconnect(); err != nil {
			return err
		}
		c.mu.RLock()
		channel = c.channel
		c.mu.RUnlock()
	}

	// channel.Consume - register consumer untuk menerima message
	// Parameters: queue, consumer-tag, auto-ack, exclusive, no-local, no-wait, arguments
	msgs, err := channel.Consume(
		c.config.Queue, // nama queue untuk consume
		"",             // consumer tag (empty = auto-generate)
		false,          // auto-ack - manual acknowledgment
		false,          // exclusive - tidak eksklusif ke consumer ini
		false,          // no-local - terima message yang dikirim connection ini
		false,          // no-wait - tunggu konfirmasi server
		nil,            // arguments tambahan
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	c.logger.WithContext("rabbitmq").Info("Started consuming messages")

	// Loop untuk menerima message dari channel
	for {
		select {
		case <-ctx.Done():
			c.logger.WithContext("rabbitmq").Info("Consumer stopped due to context cancellation")
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				c.logger.WithContext("rabbitmq").Warn("Message channel closed, attempting to reconnect")
				if err := c.reconnect(); err != nil {
					return err
				}
				return c.Consume(ctx, handler)
			}

			if err := c.processMessage(msg, handler); err != nil {
				c.logger.WithContext("rabbitmq").WithError(err).Error("Failed to process message")
			}
		}
	}
}

// processMessage memproses message yang diterima dengan error handling
// msg.Ack() - acknowledge message berhasil diproses
// msg.Nack() - negative acknowledge, message tidak berhasil diproses
func (c *Connection) processMessage(msg amqp.Delivery, handler func([]byte) error) error {
	startTime := time.Now()
	defer func() {
		if r := recover(); r != nil {
			c.metrics.IncrementFailed()
			c.logger.WithContext("rabbitmq").WithFields(map[string]interface{}{
				"panic": r,
			}).Error("Recovered from panic while processing message")
			// msg.Nack - negative acknowledgment dengan requeue
			// Parameters: multiple, requeue
			msg.Nack(false, true) // false=single message, true=requeue
		}
	}()

	if err := handler(msg.Body); err != nil {
		c.metrics.IncrementFailed()
		c.logger.WithContext("rabbitmq").WithError(err).Error("Handler returned error")
		// Nack message dan requeue untuk retry
		c.metrics.IncrementRequeued()
		return msg.Nack(false, true)
	}

	c.metrics.IncrementConsumed()
	c.metrics.RecordProcessingTime(time.Since(startTime))
	// msg.Ack - acknowledge message berhasil diproses
	// Parameters: multiple (false = acknowledge single message)
	return msg.Ack(false)
}

// Close menutup channel dan connection
// channel.Close() - tutup channel
// conn.Close() - tutup connection
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	// Tutup channel terlebih dahulu
	if c.channel != nil {
		c.channel.Close()
	}

	// Kemudian tutup connection
	if c.conn != nil {
		c.conn.Close()
	}

	c.logger.WithContext("rabbitmq").Info("RabbitMQ connection closed")
	return nil
}

// IsConnected mengecek status koneksi
// conn.IsClosed() - cek apakah connection tertutup
func (c *Connection) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return c.conn != nil && !c.conn.IsClosed() && c.channel != nil && !c.closed
}

// GetMetrics return metrics instance untuk monitoring
func (c *Connection) GetMetrics() *Metrics {
	return c.metrics
}