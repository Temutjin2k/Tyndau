
package usecase_test

import (
    "context"
    "testing"

    "github.com/Temutjin2k/Tyndau/notification-service/internal/model"
    "github.com/Temutjin2k/Tyndau/notification-service/internal/usecase"
    "github.com/Temutjin2k/Tyndau/notification-service/pkg/logger"
)

// MockEmailSender - мок для EmailSenderUseCase
type MockEmailSender struct {
    userRegisteredCalled bool
    albumReleasedCalled  bool
    email                string
    name                 string
    albumName            string
    artistName           string
}

func (m *MockEmailSender) SendUserRegisteredEmail(ctx context.Context, email, name string) error {
    m.userRegisteredCalled = true
    m.email = email
    m.name = name
    return nil
}

func (m *MockEmailSender) SendAlbumReleasedEmail(ctx context.Context, email, albumName, artistName string) error {
    m.albumReleasedCalled = true
    m.email = email
    m.albumName = albumName
    m.artistName = artistName
    return nil
}

func TestProcessEvent(t *testing.T) {
    // Создаем мок для EmailSenderUseCase
    mockEmailSender := &MockEmailSender{}
    
    // Создаем логгер
    logger := logger.NewLogger()
    
    // Создаем EventProcessorUseCase
    processor := usecase.NewEventProcessorUseCase(mockEmailSender, logger)
    
    // Тестируем обработку события user.registered
    t.Run("ProcessUserRegistered", func(t *testing.T) {
        // Создаем событие
        event := &model.Event{
            Type:   model.EventTypeUserRegistered,
            UserID: "user123",
            Email:  "user@example.com",
            Data: map[string]interface{}{
                "name": "John Doe",
            },
        }
        
        // Обрабатываем событие
        err := processor.ProcessEvent(context.Background(), event)
        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
        
        // Проверяем, что был вызван правильный метод
        if !mockEmailSender.userRegisteredCalled {
            t.Error("Expected SendUserRegisteredEmail to be called")
        }
        
        // Проверяем параметры
        if mockEmailSender.email != "user@example.com" {
            t.Errorf("Expected email to be user@example.com, got %s", mockEmailSender.email)
        }
        
        if mockEmailSender.name != "John Doe" {
            t.Errorf("Expected name to be John Doe, got %s", mockEmailSender.name)
        }
    })
    
    // Тестируем обработку события music.album_released
    t.Run("ProcessAlbumReleased", func(t *testing.T) {
        // Сбрасываем мок
        mockEmailSender.userRegisteredCalled = false
        mockEmailSender.albumReleasedCalled = false
        
        // Создаем событие
        event := &model.Event{
            Type:   model.EventTypeAlbumReleased,
            UserID: "user123",
            Email:  "user@example.com",
            Data: map[string]interface{}{
                "album_name":  "Midnight Dreams",
                "artist_name": "The Artists",
            },
        }
        
        // Обрабатываем событие
        err := processor.ProcessEvent(context.Background(), event)
        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }
        
        // Проверяем, что был вызван правильный метод
        if !mockEmailSender.albumReleasedCalled {
            t.Error("Expected SendAlbumReleasedEmail to be called")
        }
        
        // Проверяем параметры
        if mockEmailSender.email != "user@example.com" {
            t.Errorf("Expected email to be user@example.com, got %s", mockEmailSender.email)
        }
        
        if mockEmailSender.albumName != "Midnight Dreams" {
            t.Errorf("Expected album name to be Midnight Dreams, got %s", mockEmailSender.albumName)
        }
        
        if mockEmailSender.artistName != "The Artists" {
            t.Errorf("Expected artist name to be The Artists, got %s", mockEmailSender.artistName)
        }
    })
}