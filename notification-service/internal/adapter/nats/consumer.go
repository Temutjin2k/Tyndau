package nats

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/nats-io/nats.go"
    "github.com/Temutjin2k/Tyndau/notification-service/internal/model"
    "github.com/Temutjin2k/Tyndau/notification-service/pkg/logger"
)

// NatsConsumer реализует интерфейс EventConsumer для NATS
type NatsConsumer struct {
    conn       *nats.Conn
    js         nats.JetStreamContext
    streamName string
    logger     *logger.Logger
}

// NewNatsConsumer создает новый NATS consumer
func NewNatsConsumer(url, streamName string, logger *logger.Logger) (*NatsConsumer, error) {
    // Подключение к NATS
    nc, err := nats.Connect(url)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to NATS: %w", err)
    }

    // Создание JetStream контекста
    js, err := nc.JetStream()
    if err != nil {
        nc.Close()
        return nil, fmt.Errorf("failed to create JetStream context: %w", err)
    }

    // Создание стрима, если он не существует
    _, err = js.StreamInfo(streamName)
    if err != nil {
        // Стрим не существует, создаем его
        _, err = js.AddStream(&nats.StreamConfig{
            Name:     streamName,
            Subjects: []string{"user.*", "music.*"},
            MaxAge:   24 * time.Hour,
        })
        if err != nil {
            nc.Close()
            return nil, fmt.Errorf("failed to create stream: %w", err)
        }
    }

    return &NatsConsumer{
        conn:       nc,
        js:         js,
        streamName: streamName,
        logger:     logger,
    }, nil
}

// Subscribe подписывается на указанные темы
func (c *NatsConsumer) Subscribe(ctx context.Context, subjects []string, handler func(event *model.Event) error) error {
    for _, subject := range subjects {
        // Создаем durable consumer для каждой темы
        consumerName := fmt.Sprintf("notification-service-%s", subject)
        
        // Проверяем, существует ли consumer
        _, err := c.js.ConsumerInfo(c.streamName, consumerName)
        if err != nil {
            // Consumer не существует, создаем его
            _, err = c.js.AddConsumer(c.streamName, &nats.ConsumerConfig{
                Durable:       consumerName,
                AckPolicy:     nats.AckExplicitPolicy,
                FilterSubject: subject,
            })
            if err != nil {
                return fmt.Errorf("failed to create consumer for %s: %w", subject, err)
            }
        }
        
        // Подписываемся на consumer
        sub, err := c.js.PullSubscribe(subject, consumerName)
        if err != nil {
            return fmt.Errorf("failed to subscribe to %s: %w", subject, err)
        }
        
        // Запускаем горутину для обработки сообщений
        go func(subject string, subscription *nats.Subscription) {
            c.logger.Info("Started listening for events on subject: %s", subject)
            
            for {
                select {
                case <-ctx.Done():
                    c.logger.Info("Stopping subscription to %s", subject)
                    subscription.Unsubscribe()
                    return
                default:
                    // Получаем сообщения
                    msgs, err := subscription.Fetch(1, nats.MaxWait(1*time.Second))
                    if err != nil {
                        if err == nats.ErrTimeout {
                            // Нет сообщений, продолжаем
                            continue
                        }
                        c.logger.Error("Error fetching messages: %v", err)
                        time.Sleep(1 * time.Second) // Ждем перед повторной попыткой
                        continue
                    }
                    
                    for _, msg := range msgs {
                        // Логируем получение сообщения
                        c.logger.Info("Received message on subject %s: %s", subject, string(msg.Data))
                        
                        // Парсим событие
                        var event model.Event
                        if err := json.Unmarshal(msg.Data, &event); err != nil {
                            c.logger.Error("Error unmarshalling event: %v", err)
                            msg.Nak() // Отрицательное подтверждение
                            continue
                        }
                        
                        // Устанавливаем timestamp, если он не установлен
                        if event.Timestamp.IsZero() {
                            event.Timestamp = time.Now()
                        }
                        
                        // Обрабатываем событие
                        c.logger.Info("Processing event: %s", event.Type)
                        if err := handler(&event); err != nil {
                            c.logger.Error("Error handling event: %v", err)
                            msg.Nak() // Отрицательное подтверждение
                            continue
                        }
                        
                        // Подтверждаем сообщение
                        if err := msg.Ack(); err != nil {
                            c.logger.Error("Error acknowledging message: %v", err)
                        }
                    }
                }
            }
        }(subject, sub)
    }
    
    return nil
}

// Close закрывает соединение с NATS
func (c *NatsConsumer) Close() error {
    c.conn.Close()
    return nil
}