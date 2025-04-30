package smtp

import (
	"context"
	"fmt"
	"time"

	"github.com/Temutjin2k/Tyndau/notification-service/internal/model"
	"github.com/Temutjin2k/Tyndau/notification-service/pkg/logger"
	"gopkg.in/gomail.v2"
)

// SMTPConfig содержит конфигурацию SMTP
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// SMTPSender реализует интерфейс EmailSender для SMTP
type SMTPSender struct {
	config SMTPConfig
	logger *logger.Logger
}

// NewSMTPSender создает новый SMTP sender
func NewSMTPSender(config SMTPConfig, logger *logger.Logger) *SMTPSender {
	return &SMTPSender{
		config: config,
		logger: logger,
	}
}

// SendEmail отправляет email через gomail
func (s *SMTPSender) SendEmail(ctx context.Context, email *model.Email) error {
	// Таймаут
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	s.logger.Info("Sending email to %s with subject: %s", email.To, email.Subject)

	// Создаем сообщение
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/html", email.Body)

	// Настраиваем SMTP отправку
	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)

	// Отправляем в отдельной горутине
	done := make(chan error, 1)
	go func() {
		done <- d.DialAndSend(m)
	}()

	select {
	case <-ctx.Done():
		s.logger.Error("Email sending timed out: %v", ctx.Err())
		return ctx.Err()
	case err := <-done:
		if err != nil {
			s.logger.Error("Failed to send email: %v", err)
			return fmt.Errorf("failed to send email: %w", err)
		}
		s.logger.Info("Email sent successfully to %s", email.To)
		return nil
	}
}
