package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config содержит конфигурацию приложения
type Config struct {
	// NATS configuration
	NatsURL    string
	NatsStream string

	// SMTP configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	// Templates configuration
	TemplatesDir string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() (*Config, error) {
	config := &Config{}

	// NATS configuration
	config.NatsURL = getEnvString("NATS_URL", "nats://localhost:4222")
	config.NatsStream = getEnvString("NATS_STREAM", "tyndau")

	// SMTP configuration
	config.SMTPHost = getEnvString("SMTP_HOST", "smtp.mail.me.com")
	smtpPort, err := getEnvInt("SMTP_PORT", 587)
	if err != nil {
		return nil, err
	}
	config.SMTPPort = smtpPort
	config.SMTPUsername = getEnvString("SMTP_USERNAME", "")
	config.SMTPPassword = getEnvString("SMTP_PASSWORD", "")
	config.SMTPFrom = getEnvString("SMTP_FROM", "noreply@tyndau.com")

	// Templates configuration
	config.TemplatesDir = getEnvString("TEMPLATES_DIR", "./templates")

	return config, nil
}

// getEnvString получает строковую переменную окружения с значением по умолчанию
func getEnvString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// getEnvInt получает целочисленную переменную окружения с значением по умолчанию
func getEnvInt(key string, defaultValue int) (int, error) {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %w", key, err)
	}

	return value, nil
}
