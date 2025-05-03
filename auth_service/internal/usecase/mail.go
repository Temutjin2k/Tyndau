package usecase

import (
	"context"
	"errors"
	"time"
)

type MailProvider struct {
	producer Producer
}

func NewMail(producer Producer) *MailProvider {
	return &MailProvider{
		producer: producer,
	}
}

// SendWelcome sends mail and name to MessageQueue
func (m *MailProvider) SendWelcome(ctx context.Context, email, name string) error {
	if email == "" || name == "" {
		return errors.New("email and name are required")
	}

	event := map[string]interface{}{
		"event_type": "user.registered",
		"user_id":    "", // can be filled later
		"email":      email,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"data": map[string]any{
			"name": name,
		},
	}

	subject := "tyndau.user_registered"
	return m.producer.SendEvent(ctx, subject, event)
}
