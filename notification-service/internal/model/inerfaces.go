package model

import "context"

// EventConsumer определяет интерфейс для потребления событий
type EventConsumer interface {
	Subscribe(ctx context.Context, subjects []string, handler func(event *Event) error) error
	Close() error
}

// EmailSender определяет интерфейс для отправки email
type EmailSender interface {
	SendEmail(ctx context.Context, email *Email) error
}

// TemplateEngine определяет интерфейс для работы с шаблонами
type TemplateEngine interface {
	RenderTemplate(templateName string, data interface{}) (string, error)
}
