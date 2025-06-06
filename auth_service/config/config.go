package config

import (
	"time"

	nats "github.com/Temutjin2k/Tyndau/auth_service/internal/adapter/nats"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server             Server
		GRPCServices       GRPCServices
		NatsProducerConfig nats.ProducerConfig
		JWT                JWT

		Version string `env:"VERSION"`
	}

	// We can have multiple servers like gRPC or smth else.
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
		UserGRPCService UserGRPCService
	}

	UserGRPCService struct {
		Addr string `env:"USER_GRPC_ADDR,notEmpty"`
	}

	JWT struct {
		Secret string `env:"JWT_SECRET,notEmpty"`
	}
)

func New() (*Config, error) {
	var cfg Config

	// Load environment variables from .env file
	godotenv.Load()

	// Parse environment variables into the Config structure
	err := env.Parse(&cfg)

	PrintConfig(cfg)
	return &cfg, err
}
