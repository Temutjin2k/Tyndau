package config

import (
	"time"

	nats "github.com/Temutjin2k/Tyndau/music-service/internal/adapter/nats"
	"github.com/Temutjin2k/Tyndau/music-service/internal/adapter/storage"
	"github.com/Temutjin2k/Tyndau/music-service/pkg/postgres"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server             Server
		NatsProducerConfig nats.ProducerConfig
		Postgres           postgres.Config
		Redis              Redis
		Minio              storage.MinioConfig

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

	Redis struct {
		Addr     string        `env:"REDIS_ADDR" envDefault:"localhost:6379"`
		Password string        `env:"REDIS_PASSWORD,required"`
		DB       int           `env:"REDIS_DB" envDefault:"0"`
		TTL      time.Duration `env:"CACHE_TTL_SECONDS" envDefault:"300s"`
	}
)

func New() (*Config, error) {
	var cfg Config

	// Load environment variables from .env file
	godotenv.Load()

	// Parse environment variables into the Config structure
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	// Print config for debug
	PrintConfig(cfg)
	return &cfg, nil
}
