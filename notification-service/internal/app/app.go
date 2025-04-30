package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv" // Add this import for string to int conversion
	"syscall"

	"github.com/Temutjin2k/Tyndau/notification-service/internal/adapter/nats"
	"github.com/Temutjin2k/Tyndau/notification-service/internal/adapter/smtp"
	"github.com/Temutjin2k/Tyndau/notification-service/internal/usecase"
	"github.com/Temutjin2k/Tyndau/notification-service/pkg/logger"
)

// App represents the application
type App struct {
	natsConsumer    *nats.Consumer
	eventProcessor  *usecase.EventProcessorUseCase
	emailSender     *usecase.EmailSenderUseCase
	templateEngine  *usecase.GoTemplateEngine
	logger          *logger.Logger
}

// NewApp creates a new application instance
func NewApp() *App {
	// Initialize logger
	logger := logger.NewLogger()
	
	// Initialize template engine
	templateEngine := usecase.NewGoTemplateEngine(os.Getenv("TEMPLATES_DIR"), logger)

	// Parse SMTP port from environment variable
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpPort := 587 // Default port
	if smtpPortStr != "" {
		port, err := strconv.Atoi(smtpPortStr)
		if err != nil {
			logger.Warn("Invalid SMTP_PORT value: %s, using default 587", smtpPortStr)
		} else {
			smtpPort = port
		}
	}

	// Initialize email sender
	smtpSender := smtp.NewSMTPSender(smtp.SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     smtpPort, // Now using the parsed integer
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}, logger)
	
	emailSender := usecase.NewEmailSenderUseCase(smtpSender, templateEngine, logger)

	// Initialize event processor
	eventProcessor := usecase.NewEventProcessorUseCase(emailSender, logger)

	// Initialize NATS consumer
	natsConsumer, err := nats.NewConsumer(eventProcessor)
	if err != nil {
		log.Fatalf("Failed to create NATS consumer: %v", err)
	}

	return &App{
		natsConsumer:   natsConsumer,
		eventProcessor: eventProcessor,
		emailSender:    emailSender,
		templateEngine: templateEngine,
		logger:         logger,
	}
}

// Run starts the application
func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to events
	if err := a.natsConsumer.SubscribeToEvents(ctx); err != nil {
		return err
	}

	// Wait for termination signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	a.logger.Info("Shutting down...")
	return nil
}