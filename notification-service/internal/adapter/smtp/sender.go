package smtp

import (
    "context"
    "fmt"
    "net/smtp"
    "time"

    "github.com/Temutjin2k/Tyndau/notification-service/internal/model"
    "github.com/Temutjin2k/Tyndau/notification-service/pkg/logger"
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

// SendEmail отправляет email через SMTP
func (s *SMTPSender) SendEmail(ctx context.Context, email *model.Email) error {
    // Создаем timeout контекст
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    
    // Создаем аутентификацию
    auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
    
    // Формируем сообщение
    msg := []byte(fmt.Sprintf("From: %s\r\n"+
        "To: %s\r\n"+
        "Subject: %s\r\n"+
        "MIME-Version: 1.0\r\n"+
        "Content-Type: text/html; charset=UTF-8\r\n"+
        "\r\n"+
        "%s\r\n", s.config.From, email.To, email.Subject, email.Body))
    
    // Логируем отправку
    s.logger.Info("Sending email to %s with subject: %s", email.To, email.Subject)
    
    // Отправляем email
    done := make(chan error, 1)
    go func() {
        done <- smtp.SendMail(
            fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
            auth,
            s.config.From,
            []string{email.To},
            msg,
        )
    }()
    
    // Ждем отправки или отмены контекста
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