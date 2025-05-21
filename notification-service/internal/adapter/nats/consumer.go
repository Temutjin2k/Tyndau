package nats

import (
	"context"
	"encoding/json" // Импортируем json
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/Temutjin2k/Tyndau/notification-service/internal/model"
	"github.com/Temutjin2k/Tyndau/notification-service/pkg/logger"
)

type Consumer struct {
	js      jetstream.JetStream
	stream  string
	handler model.EventProcessor
	logger  *logger.Logger
}

func NewConsumer(handler model.EventProcessor) (*Consumer, error) {
	log.Println("Creating NATS consumer")

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	log.Printf("NATS URL: %s", natsURL)

	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	stream := os.Getenv("NATS_STREAM")
	if stream == "" {
		stream = "tyndau"
	}

	log.Printf("NATS Stream: %s", stream)

	_, err = js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name:     stream,
		Subjects: []string{fmt.Sprintf("%s.>", stream)},
	})
	if err != nil && err != jetstream.ErrStreamNameAlreadyInUse {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	return &Consumer{
		js:      js,
		stream:  stream,
		handler: handler,
	}, nil
}

func (c *Consumer) SubscribeToEvents(ctx context.Context) error {
	consumerPrefix := os.Getenv("NATS_CONSUMER_PREFIX")
	if consumerPrefix == "" {
		consumerPrefix = "tyndau_consumer"
	}

	// Используем строковые значения, а не константы
	subjects := map[string]string{
		"user_registered": "user.registered",
		"album_released":  "music.album_released",
	}

	for subject, eventType := range subjects {
		fullSubject := fmt.Sprintf("%s.%s", c.stream, subject)
		consumerName := fmt.Sprintf("%s_%s", consumerPrefix, subject)

		log.Printf("Setting up consumer for subject: %s (eventType: %s)", fullSubject, eventType)

		consumerConfig := jetstream.ConsumerConfig{
			Durable:       consumerName,
			AckPolicy:     jetstream.AckExplicitPolicy,
			MaxDeliver:    3,
			AckWait:       30 * time.Second,
			FilterSubject: fullSubject,
		}

		consumer, err := c.js.CreateOrUpdateConsumer(ctx, c.stream, consumerConfig)
		if err != nil {
			return fmt.Errorf("failed to create consumer for %s: %w", subject, err)
		}

		_, err = consumer.Consume(func(msg jetstream.Msg) {
			log.Printf("📥 Received message on subject: %s", msg.Subject())
			log.Printf("📦 Message data: %s", string(msg.Data()))

			var event model.Event
			if err := json.Unmarshal(msg.Data(), &event); err != nil {
				log.Printf("❌ Failed to unmarshal event: %v", err)
				return
			}

			var processErr error
			switch event.Type {
			case "user.registered": // Используем строку для типа события
				processErr = c.handler.ProcessUserRegistered(msg.Data())
			case "music.album_released": // Используем строку для типа события
				processErr = c.handler.ProcessAlbumReleased(msg.Data())
			default:
				log.Printf("⚠️ Unknown event type: %s", event.Type)
			}

			if processErr != nil {
				log.Printf("❌ Error processing %s: %v", event.Type, processErr)
			}

			if err := msg.Ack(); err != nil {
				log.Printf("❌ Failed to ACK message: %v", err)
			}
		})

		if err != nil {
			return fmt.Errorf("failed to consume from subject %s: %w", subject, err)
		}

		log.Printf("✅ Subscribed to %s (%s)", fullSubject, eventType)
	}

	return nil
}
