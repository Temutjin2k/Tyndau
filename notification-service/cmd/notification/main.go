package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Temutjin2k/Tyndau/notification-service/config"
	"github.com/Temutjin2k/Tyndau/notification-service/internal/adapter/nats"
	"github.com/Temutjin2k/Tyndau/notification-service/internal/adapter/smtp"
	"github.com/Temutjin2k/Tyndau/notification-service/internal/usecase"
	"github.com/Temutjin2k/Tyndau/notification-service/pkg/logger"
)

func main() {
	// Инициализация логгера
	logger := logger.NewLogger()
	logger.Info("Starting notification service")

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	// Создание контекста с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация template engine
	templateEngine := usecase.NewGoTemplateEngine(cfg.TemplatesDir, logger)

	// Инициализация SMTP sender
	smtpSender := smtp.NewSMTPSender(smtp.SMTPConfig{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
	}, logger)

	// Инициализация email sender use case
	emailSender := usecase.NewEmailSenderUseCase(smtpSender, templateEngine, logger)

	// Инициализация event processor use case
	eventProcessor := usecase.NewEventProcessorUseCase(emailSender, logger)

	// Инициализация NATS consumer - FIXED: using NewConsumer instead of NewNatsConsumer
	natsConsumer, err := nats.NewConsumer(eventProcessor)
	if err != nil {
		logger.Fatal("Failed to create NATS consumer: %v", err)
	}

	// Подписка на события
	err = natsConsumer.SubscribeToEvents(ctx)
	if err != nil {
		logger.Fatal("Failed to subscribe to events: %v", err)
	}

	logger.Info("Notification service started, listening for events")

	// Ожидание сигнала завершения
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("Received shutdown signal")

	// Отмена контекста для остановки подписок
	cancel()

	// Даем время на завершение текущих операций
	time.Sleep(1 * time.Second)

	logger.Info("Notification service shutdown complete")
}