package config

import (
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

	err := godotenv.Load()
	if err != nil {
		return &cfg, err
	}

	err = env.Parse(&cfg)

	return &cfg, err
}
