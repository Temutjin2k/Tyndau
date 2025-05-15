package config

import (
	"path/filepath"
	"time"

	nats "github.com/Temutjin2k/Tyndau/music-service/internal/adapter/nats"
	"github.com/Temutjin2k/Tyndau/music-service/pkg/postgres"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server             Server
		GRPCServices       GRPCServices
		NatsProducerConfig nats.ProducerConfig
		Postgres           postgres.Config

		Version string `env:"VERSION"`
	}

	Server struct {
		GRPCServer GRPCServer
	}

	GRPCServer struct {
		Port                  int           `env:"GRPC_PORT" envDefault:"50051"`
		MaxRecvMsgSizeMiB     int           `env:"GRPC_MAX_MESSAGE_SIZE_MIB" envDefault:"12"`
		MaxConnectionAge      time.Duration `env:"GRPC_MAX_CONNECTION_AGE" envDefault:"30s"`
		MaxConnectionAgeGrace time.Duration `env:"GRPC_MAX_CONNECTION_AGE_GRACE" envDefault:"10s"`
	}

	GRPCServices struct {
	}
)

func New() (*Config, error) {
	var cfg Config

	// Пытаемся загрузить .env
	if err := godotenv.Load(filepath.Join("cmd", "music", ".env")); err != nil {
		return &cfg, err
	}

	// Парсим переменные окружения
	if err := env.Parse(&cfg); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
