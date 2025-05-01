package model

import "context"

// EventProcessor processes events
type EventProcessor interface {
	ProcessEvent(ctx context.Context, event *Event) error
	// Add this method to match what's being called in consumer.go
	ProcessUserRegistered(data []byte) error
	ProcessAlbumReleased(data []byte) error
}

// EmailSender sends emails
type EmailSender interface {
	SendEmail(ctx context.Context, email *Email) error
}

// EmailSenderUseCase sends different types of emails
type EmailSenderUseCase interface {
	SendUserRegisteredEmail(ctx context.Context, email, name string) error
	SendAlbumReleasedEmail(ctx context.Context, email, albumName, artistName string) error
}

// TemplateEngine renders templates
type TemplateEngine interface {
	RenderTemplate(templateName string, data interface{}) (string, error)
}

// EventConsumer is the interface that consumer.go expects
type EventConsumer interface {
	ProcessUserRegistered(data []byte) error
}