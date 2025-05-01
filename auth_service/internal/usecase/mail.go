package usecase

import "context"

type welcomeMsg struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

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
	req := welcomeMsg{
		Email: email,
		Name:  name,
	}

	m.producer.PublishWithContext(ctx, req)
	return nil
}
