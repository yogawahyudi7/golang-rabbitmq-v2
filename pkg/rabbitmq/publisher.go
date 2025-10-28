package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"golang-rabbitmq-v2/pkg/logger"
)

type Publisher interface {
	Publish(ctx context.Context, routingKey string, body []byte) error
	PublishWithConfirm(ctx context.Context, routingKey string, body []byte) error
	Close() error
}

type publisher struct {
	conn         *Connection
	logger       *logger.Logger
	confirms     chan amqp.Confirmation
	confirmMutex sync.RWMutex
	closed       bool
}

// NewPublisher membuat publisher baru dengan publisher confirms
// Menginisialisasi confirmation channel dan setup confirms mode
func NewPublisher(conn *Connection, log *logger.Logger) (Publisher, error) {
	p := &publisher{
		conn:     conn,
		logger:   log,
		confirms: make(chan amqp.Confirmation, 1000), // buffer untuk confirmations
	}

	if err := p.setupConfirms(); err != nil {
		return nil, fmt.Errorf("failed to setup publisher confirms: %w", err)
	}

	return p, nil
}

// setupConfirms mengaktifkan publisher confirmation mode
// channel.Confirm() - aktifkan confirm mode pada channel
// channel.NotifyPublish() - register channel untuk menerima confirmations
func (p *publisher) setupConfirms() error {
	p.conn.mu.RLock()
	channel := p.conn.channel
	p.conn.mu.RUnlock()

	if channel == nil {
		return fmt.Errorf("channel is nil")
	}

	// channel.Confirm - aktifkan publisher confirmation mode
	// Parameter: no-wait (false = tunggu konfirmasi server)
	if err := channel.Confirm(false); err != nil {
		return fmt.Errorf("failed to put channel in confirm mode: %w", err)
	}

	// channel.NotifyPublish - register channel untuk menerima publish confirmations
	// Return channel yang akan menerima amqp.Confirmation untuk setiap publish
	notifyConfirm := channel.NotifyPublish(make(chan amqp.Confirmation, 1000))
	
	// Goroutine untuk handle confirmations secara asynchronous
	go p.handleConfirmations(notifyConfirm)

	p.logger.WithContext("rabbitmq-publisher").Info("Publisher confirms enabled")
	return nil
}

// handleConfirmations memproses confirmation yang diterima dari broker
// Menerima amqp.Confirmation dan meneruskan ke internal confirmation channel
func (p *publisher) handleConfirmations(confirms <-chan amqp.Confirmation) {
	for confirm := range confirms {
		select {
		case p.confirms <- confirm:
			// Confirmation berhasil diteruskan ke internal channel
		default:
			// Channel penuh, drop confirmation (tidak ideal dalam production)
			p.logger.WithContext("rabbitmq-publisher").Warn("Confirmation channel full, dropping confirmation")
		}
	}
}

// Publish mengirim message tanpa menunggu confirmation (fire-and-forget)
// Lebih cepat tapi tidak ada guarantee message sampai ke broker
func (p *publisher) Publish(ctx context.Context, routingKey string, body []byte) error {
	p.confirmMutex.RLock()
	if p.closed {
		p.confirmMutex.RUnlock()
		return fmt.Errorf("publisher is closed")
	}
	p.confirmMutex.RUnlock()

	// Track publishing start
	startTime := time.Now()
	p.conn.metrics.StartPublishing()
	defer p.conn.metrics.EndPublishing()

	// Delegasi ke method publish pada Connection (tanpa wait confirmation)
	err := p.conn.Publish(ctx, routingKey, body)

	// Track event
	duration := time.Since(startTime)
	p.conn.metrics.TrackPublishEvent(routingKey, err == nil, duration, err)

	return err
}

// PublishWithConfirm mengirim message dan menunggu confirmation dari broker
// Memberikan guarantee bahwa message sudah diterima oleh broker
func (p *publisher) PublishWithConfirm(ctx context.Context, routingKey string, body []byte) error {
	p.confirmMutex.RLock()
	if p.closed {
		p.confirmMutex.RUnlock()
		return fmt.Errorf("publisher is closed")
	}
	p.confirmMutex.RUnlock()

	// Track publishing start
	startTime := time.Now()
	p.conn.metrics.StartPublishing()
	defer p.conn.metrics.EndPublishing()

	// Publish message ke broker
	if err := p.conn.Publish(ctx, routingKey, body); err != nil {
		duration := time.Since(startTime)
		p.conn.metrics.TrackPublishEvent(routingKey, false, duration, err)
		return err
	}

	// Tunggu confirmation dari broker
	err := p.waitForConfirmation(ctx)

	// Track event
	duration := time.Since(startTime)
	p.conn.metrics.TrackPublishEvent(routingKey, err == nil, duration, err)

	return err
}

// waitForConfirmation menunggu confirmation dari broker dengan timeout
// Menerima amqp.Confirmation yang berisi Ack/Nack status dari broker
func (p *publisher) waitForConfirmation(ctx context.Context) error {
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case confirm := <-p.confirms:
		// Cek apakah message di-acknowledge atau di-negative-acknowledge
		if confirm.Ack {
			p.conn.metrics.IncrementConfirmAcked()
			p.conn.metrics.TrackConfirmEvent(confirm.DeliveryTag, true)
			p.logger.WithContext("rabbitmq-publisher").WithFields(map[string]interface{}{
				"delivery_tag": confirm.DeliveryTag, // unique ID untuk message
			}).Debug("Message confirmed")
			return nil
		}
		// Message di-nack oleh broker (misalnya: disk penuh, queue tidak ada)
		p.conn.metrics.IncrementConfirmNacked()
		p.conn.metrics.TrackConfirmEvent(confirm.DeliveryTag, false)
		return fmt.Errorf("message nacked by broker, delivery tag: %d", confirm.DeliveryTag)

	case <-timer.C:
		// Timeout menunggu confirmation
		p.conn.metrics.IncrementConfirmTimeout()
		return fmt.Errorf("timeout waiting for confirmation after %s", timeout)

	case <-ctx.Done():
		// Context cancelled
		return ctx.Err()
	}
}

// Close menutup publisher dan cleanup resources
// Menutup confirmation channel dan menandai publisher sebagai closed
func (p *publisher) Close() error {
	p.confirmMutex.Lock()
	defer p.confirmMutex.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true
	close(p.confirms) // Tutup confirmation channel

	p.logger.WithContext("rabbitmq-publisher").Info("Publisher closed")
	return nil
}