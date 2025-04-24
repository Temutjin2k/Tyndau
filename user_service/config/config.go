package config

import (
	"time"
	"user_service/pkg/postgres"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Postgres postgres.Config
		Server   Server

		Version string `env:"VERSION"`
	}

	// We can have multiple servers like gRPC or smth else.
	Server struct {
		HTTPServer HTTPServer
		GRPCServer GRPCServer
	}

	HTTPServer struct {
		Port           int           `env:"HTTP_PORT" envDefault:"8082"`
		ReadTimeout    time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"30s"`
		WriteTimeout   time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"30s"`
		IdleTimeout    time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
		MaxHeaderBytes int           `env:"HTTP_MAX_HEADER_BYTES" envDefault:"1048576"` // 1 MB
		TrustedProxies []string      `env:"HTTP_TRUSTED_PROXIES" envSeparator:","`
		Mode           string        `env:"GIN_MODE" envDefault:"release"` // Can be: release, debug, test
	}

	GRPCServer struct {
		Port                  int16         `env:"GRPC_PORT,notEmpty"`
		MaxRecvMsgSizeMiB     int           `env:"GRPC_MAX_MESSAGE_SIZE_MIB" envDefault:"12"`
		MaxConnectionAge      time.Duration `env:"GRPC_MAX_CONNECTION_AGE" envDefault:"30s"`
		MaxConnectionAgeGrace time.Duration `env:"GRPC_MAX_CONNECTION_AGE_GRACE" envDefault:"10s"`
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
