package usecase

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/Temutjin2k/Tyndau/notification-service/internal/model"
    "github.com/Temutjin2k/Tyndau/notification-service/pkg/logger"
)

// EventProcessorUseCase обрабатывает события
type EventProcessorUseCase struct {
    emailSender *EmailSenderUseCase
    logger      *logger.Logger
}

// NewEventProcessorUseCase создает новый use case для обработки событий
func NewEventProcessorUseCase(
    emailSender *EmailSenderUseCase,
    logger *logger.Logger,
) *EventProcessorUseCase {
    return &EventProcessorUseCase{
        emailSender: emailSender,
        logger:      logger,
    }
}

// ProcessEvent обрабатывает событие
func (uc *EventProcessorUseCase) ProcessEvent(ctx context.Context, event *model.Event) error {
    uc.logger.Info("Processing event: %s for user %s", event.Type, event.UserID)
    
    switch event.Type {
    case model.EventTypeUserRegistered:
        return uc.processUserRegistered(ctx, event)
    case model.EventTypeAlbumReleased:
        return uc.processAlbumReleased(ctx, event)
    default:
        uc.logger.Warn("Unknown event type: %s", event.Type)
        return fmt.Errorf("unknown event type: %s", event.Type)
    }
}

// processUserRegistered обрабатывает событие регистрации пользователя
func (uc *EventProcessorUseCase) processUserRegistered(ctx context.Context, event *model.Event) error {
    // Парсим данные события
    var data model.UserRegisteredData
    if err := parseEventData(event.Data, &data); err != nil {
        return fmt.Errorf("failed to parse user registered data: %w", err)
    }
    
    // Отправляем email
    return uc.emailSender.SendUserRegisteredEmail(ctx, event.Email, data.Name)
}

// processAlbumReleased обрабатывает событие релиза альбома
func (uc *EventProcessorUseCase) processAlbumReleased(ctx context.Context, event *model.Event) error {
    // Парсим данные события
    var data model.AlbumReleasedData
    if err := parseEventData(event.Data, &data); err != nil {
        return fmt.Errorf("failed to parse album released data: %w", err)
    }
    
    // Отправляем email
    return uc.emailSender.SendAlbumReleasedEmail(ctx, event.Email, data.AlbumName, data.ArtistName)
}

// parseEventData парсит данные события
func parseEventData(data interface{}, target interface{}) error {
    // Если данные уже в нужном формате, просто копируем
    if value, ok := data.(map[string]interface{}); ok {
        // Преобразуем map в JSON и обратно в структуру
        jsonData, err := json.Marshal(value)
        if err != nil {
            return err
        }
        return json.Unmarshal(jsonData, target)
    }
    
    // Если данные в виде JSON строки
    if jsonStr, ok := data.(string); ok {
        return json.Unmarshal([]byte(jsonStr), target)
    }
    
    return fmt.Errorf("unsupported data format")
}