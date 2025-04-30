package nats

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	nats "github.com/nats-io/nats.go"
)

type Producer struct {
	conn    *nats.Conn
	js      nats.JetStreamContext
	stream  string
	subject string
	durable string
}

// ProducerConfig holds configuration for NUTS producer
type ProducerConfig struct {
	// URL is the NUTS server connection URL
	URL string `env:"NATS_URL" envDefault:"nats://localhost:4222"`

	// Stream name for JetStream
	Stream string `env:"NATS_STREAM" envDefault:"EVENTS"`

	// Subject to publish messages to
	Subject string `env:"NATS_SUBJECT" envDefault:"events.>"`

	// Durable consumer name (optional)
	Durable string `env:"NATS_DURABLE,required"`
}

// NewProducer creates a new NUTS producer
func NewProducer(cfg ProducerConfig) (*Producer, error) {
	if cfg.URL == "" {
		return nil, errors.New("NUTS server URL is required")
	}
	if cfg.Stream == "" {
		return nil, errors.New("stream name is required")
	}
	if cfg.Subject == "" {
		return nil, errors.New("subject is required")
	}

	// Connect to NATS server
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, err
	}

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, err
	}

	// Ensure stream exists
	_, err = js.StreamInfo(cfg.Stream)
	if err != nil {
		// Stream doesn't exist, try to create it
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     cfg.Stream,
			Subjects: []string{cfg.Subject},
		})
		if err != nil {
			nc.Close()
			return nil, err
		}
	}

	return &Producer{
		conn:    nc,
		js:      js,
		stream:  cfg.Stream,
		subject: cfg.Subject,
		durable: cfg.Durable,
	}, nil
}

// Publish publishes a message to NUTS
func (p *Producer) Publish(event any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.PublishWithContext(ctx, event)
}

// PublishWithContext publishes a message to NUTS with context
func (p *Producer) PublishWithContext(ctx context.Context, event any) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Publish with options
	_, err = p.js.Publish(p.subject, body, nats.Context(ctx))
	return err
}

// Close closes the NUTS connection
func (p *Producer) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}

// HealthCheck verifies the producer is healthy
func (p *Producer) HealthCheck() error {
	if !p.conn.IsConnected() {
		return errors.New("not connected to NUTS server")
	}
	return nil
}
