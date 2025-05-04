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

// SendBulkEmail отправляет email с использованием BCC для массовой рассылки
func (s *SMTPSender) SendBulkEmail(ctx context.Context, emails []string, subject, body string, batchSize int) error {
	if len(emails) == 0 {
		return nil
	}

	// Если batchSize не указан или некорректный, используем разумное значение по умолчанию
	if batchSize <= 0 {
		batchSize = 50
	}

	// Разбиваем список получателей на партии
	var batches [][]string
	for i := 0; i < len(emails); i += batchSize {
		end := i + batchSize
		if end > len(emails) {
			end = len(emails)
		}
		batches = append(batches, emails[i:end])
	}

	s.logger.Info("Sending bulk email to %d recipients in %d batches", len(emails), len(batches))

	// Отправляем каждую партию
	for i, batch := range batches {
		// Таймаут для каждой партии
		batchCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		// Создаем сообщение
		m := gomail.NewMessage()
		m.SetHeader("From", s.config.From)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)

		// Добавляем первого получателя в To (чтобы письмо было валидным)
		if len(batch) > 0 {
			m.SetHeader("To", s.config.From) // Отправляем себе, чтобы не раскрывать получателей
		}

		// Добавляем получателей в BCC
		for _, email := range batch {
			m.SetHeader("Bcc", email) // Правильный способ добавления BCC в gomail
		}

		// Настраиваем SMTP отправку
		d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)

		// Отправляем в отдельной горутине
		done := make(chan error, 1)
		go func() {
			done <- d.DialAndSend(m)
		}()

		select {
		case <-batchCtx.Done():
			s.logger.Error("Batch %d/%d email sending timed out: %v", i+1, len(batches), batchCtx.Err())
			return batchCtx.Err()
		case err := <-done:
			if err != nil {
				s.logger.Error("Failed to send batch %d/%d: %v", i+1, len(batches), err)
				return fmt.Errorf("failed to send email batch %d/%d: %w", i+1, len(batches), err)
			}
			s.logger.Info("Successfully sent batch %d/%d to %d recipients", i+1, len(batches), len(batch))
		}

		// Добавляем небольшую задержку между партиями, чтобы не перегружать SMTP сервер
		if i < len(batches)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	s.logger.Info("Bulk email sent successfully to %d recipients", len(emails))
	return nil
}
