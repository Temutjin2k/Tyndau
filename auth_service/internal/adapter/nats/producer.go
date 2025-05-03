package nats

import (
	"context"
	"encoding/json"
	"errors"

	nats "github.com/nats-io/nats.go"
)

type ProducerConfig struct {
	URL string `env:"NATS_URL"`
}

type Producer struct {
	conn *nats.Conn
	cfg  ProducerConfig
}

// NewProducer connects to a NATS server using plain publish-subscribe (no JetStream)
func NewProducer(cfg ProducerConfig) (*Producer, error) {
	if cfg.URL == "" {
		return nil, errors.New("NATS URL is required")
	}

	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, err
	}

	return &Producer{conn: nc, cfg: cfg}, nil
}

// SendEvent publishes a JSON message to a subject using the basic NATS Publish
func (p *Producer) SendEvent(ctx context.Context, subject string, event map[string]any) error {
	if subject == "" {
		return errors.New("subject cannot be empty")
	}
	if event == nil {
		return errors.New("event data cannot be nil")
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Create a done channel to wait for the publish or context timeout
	done := make(chan error, 1)

	go func() {
		err := p.conn.Publish(subject, body)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// Close closes the NATS connection
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
