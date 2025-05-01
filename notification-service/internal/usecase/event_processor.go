package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Temutjin2k/Tyndau/notification-service/internal/model"
	"github.com/Temutjin2k/Tyndau/notification-service/pkg/logger"
)

type EventProcessorUseCase struct {
	emailSender model.EmailSenderUseCase
	logger      *logger.Logger
}

func NewEventProcessorUseCase(emailSender model.EmailSenderUseCase, logger *logger.Logger) *EventProcessorUseCase {
	return &EventProcessorUseCase{
		emailSender: emailSender,
		logger:      logger,
	}
}

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

func (uc *EventProcessorUseCase) ProcessUserRegistered(data []byte) error {
	var event model.Event
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	var userData model.UserRegisteredData
	if err := parseEventData(event.Data, &userData); err != nil {
		return fmt.Errorf("failed to parse user registered data: %w", err)
	}

	log.Printf("📩 Welcome email sent to %s (%s)", userData.Name, event.Email)
	return uc.emailSender.SendUserRegisteredEmail(context.Background(), event.Email, userData.Name)
}

func (uc *EventProcessorUseCase) ProcessAlbumReleased(data []byte) error {
	var event model.Event
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	var albumData model.AlbumReleasedData
	if err := parseEventData(event.Data, &albumData); err != nil {
		return fmt.Errorf("failed to parse album released data: %w", err)
	}

	log.Printf("🎶 Album released: \"%s\" by %s at %s", albumData.AlbumName, albumData.ArtistName, event.Timestamp.Format(time.RFC1123))
	return uc.emailSender.SendAlbumReleasedEmail(context.Background(), event.Email, albumData.AlbumName, albumData.ArtistName)
}

func (uc *EventProcessorUseCase) processUserRegistered(ctx context.Context, event *model.Event) error {
	var data model.UserRegisteredData
	if err := parseEventData(event.Data, &data); err != nil {
		return fmt.Errorf("failed to parse user registered data: %w", err)
	}
	return uc.emailSender.SendUserRegisteredEmail(ctx, event.Email, data.Name)
}

func (uc *EventProcessorUseCase) processAlbumReleased(ctx context.Context, event *model.Event) error {
	var data model.AlbumReleasedData
	if err := parseEventData(event.Data, &data); err != nil {
		return fmt.Errorf("failed to parse album released data: %w", err)
	}
	return uc.emailSender.SendAlbumReleasedEmail(ctx, event.Email, data.AlbumName, data.ArtistName)
}

func parseEventData(data interface{}, target interface{}) error {
	switch v := data.(type) {
	case map[string]interface{}:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return json.Unmarshal(jsonData, target)
	case string:
		return json.Unmarshal([]byte(v), target)
	default:
		return fmt.Errorf("unsupported data format: %T", v)
	}
}
