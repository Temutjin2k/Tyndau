package usecase

import (
	"context"
	"fmt"
	"sync"

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

// SendAlbumReleasedMassEmail отправляет уведомление о релизе альбома на множество адресов
func (uc *EmailSenderUseCase) SendAlbumReleasedMassEmail(ctx context.Context, emails []string, albumName, artistName string) error {
	if len(emails) == 0 {
		uc.logger.Warn("No emails provided for mass album release notification")
		return nil
	}

	// Данные для шаблона
	data := map[string]interface{}{
		"AlbumName":  albumName,
		"ArtistName": artistName,
	}

	// Рендерим шаблон (делаем это один раз для всех писем)
	body, err := uc.templateEngine.RenderTemplate("album_released", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Создаем базовый объект email
	subject := fmt.Sprintf("New Album Release: %s by %s", albumName, artistName)

	// Используем WaitGroup для отслеживания всех горутин
	var wg sync.WaitGroup
	errChan := make(chan error, len(emails))

	// Отправляем email на каждый адрес параллельно
	for _, email := range emails {
		wg.Add(1)
		go func(emailAddr string) {
			defer wg.Done()

			// Создаем email для конкретного получателя
			emailObj := &model.Email{
				To:      emailAddr,
				Subject: subject,
				Body:    body,
			}

			// Отправляем email
			if err := uc.emailSender.SendEmail(ctx, emailObj); err != nil {
				uc.logger.Error("Failed to send email to %s: %v", emailAddr, err)
				errChan <- fmt.Errorf("failed to send email to %s: %w", emailAddr, err)
			} else {
				uc.logger.Info("Successfully sent album release email to %s", emailAddr)
			}
		}(email)
	}

	// Ждем завершения всех горутин
	wg.Wait()
	close(errChan)

	// Проверяем, были ли ошибки
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to send %d out of %d emails", len(errs), len(emails))
	}

	return nil
}