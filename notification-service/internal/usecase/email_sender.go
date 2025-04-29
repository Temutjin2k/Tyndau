package usecase

import (
	"context"
	"fmt"

	"github.com/Temutjin2k/Tyndau/notification-service/internal/model"
	"github.com/Temutjin2k/Tyndau/notification-service/pkg/logger"
)

// EmailSenderUseCase реализует логику отправки email
type EmailSenderUseCase struct {
	emailSender    model.EmailSender
	templateEngine model.TemplateEngine
	logger         *logger.Logger
}

// NewEmailSenderUseCase создает новый use case для отправки email
func NewEmailSenderUseCase(
	emailSender model.EmailSender,
	templateEngine model.TemplateEngine,
	logger *logger.Logger,
) *EmailSenderUseCase {
	return &EmailSenderUseCase{
		emailSender:    emailSender,
		templateEngine: templateEngine,
		logger:         logger,
	}
}

// SendUserRegisteredEmail отправляет приветственное письмо
func (uc *EmailSenderUseCase) SendUserRegisteredEmail(ctx context.Context, email string, name string) error {
	// Данные для шаблона
	data := map[string]interface{}{
		"Name": name,
	}

	// Рендерим шаблон
	body, err := uc.templateEngine.RenderTemplate("user_registered", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Создаем email
	emailObj := &model.Email{
		To:      email,
		Subject: "Welcome to our service!",
		Body:    body,
	}

	// Отправляем email
	return uc.emailSender.SendEmail(ctx, emailObj)
}

// SendAlbumReleasedEmail отправляет уведомление о релизе альбома
func (uc *EmailSenderUseCase) SendAlbumReleasedEmail(ctx context.Context, email string, albumName, artistName string) error {
	// Данные для шаблона
	data := map[string]interface{}{
		"AlbumName":  albumName,
		"ArtistName": artistName,
	}

	// Рендерим шаблон
	body, err := uc.templateEngine.RenderTemplate("album_released", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Создаем email
	emailObj := &model.Email{
		To:      email,
		Subject: fmt.Sprintf("New Album Release: %s by %s", albumName, artistName),
		Body:    body,
	}

	// Отправляем email
	return uc.emailSender.SendEmail(ctx, emailObj)
}
