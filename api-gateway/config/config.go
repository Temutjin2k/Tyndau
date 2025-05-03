package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type (
	// Config holds the entire application configuration
	Config struct {
		Server  Server
		Version string `env:"VERSION"`
	}

	// Server holds the configuration for both HTTP and gRPC servers
	Server struct {
		HTTPServer      HTTPServer
		UserGRPCServers UserGRPCServer
		AuthGRPCServer  AuthGRPCServer
		MusicGRPCServer MusicGRPCServer
	}

	// HTTPServer holds the configuration for the HTTP server
	HTTPServer struct {
		Port           int           `env:"HTTP_PORT" envDefault:"8080"`
		ReadTimeout    time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"30s"`
		WriteTimeout   time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"30s"`
		IdleTimeout    time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
		MaxHeaderBytes int           `env:"HTTP_MAX_HEADER_BYTES" envDefault:"1048576"` // 1 MB
		TrustedProxies []string      `env:"HTTP_TRUSTED_PROXIES" envSeparator:","`
	}

	UserGRPCServer struct {
		Addr         string        `env:"USER_GRPC_ADDRESS,notEmpty"` // Default port for gRPC server
		ReadTimeout  time.Duration `env:"GRPC_READ_TIMEOUT" envDefault:"30s"`
		WriteTimeout time.Duration `env:"GRPC_WRITE_TIMEOUT" envDefault:"30s"`
	}

	AuthGRPCServer struct {
		Addr         string        `env:"AUTH_GRPC_ADDRESS,notEmpty"` // Default port for gRPC server
		ReadTimeout  time.Duration `env:"GRPC_READ_TIMEOUT" envDefault:"30s"`
		WriteTimeout time.Duration `env:"GRPC_WRITE_TIMEOUT" envDefault:"30s"`
	}

	MusicGRPCServer struct {
		Addr         string        `env:"MUSIC_GRPC_ADDRESS,notEmpty"` // Default port for gRPC server
		ReadTimeout  time.Duration `env:"GRPC_READ_TIMEOUT" envDefault:"30s"`
		WriteTimeout time.Duration `env:"GRPC_WRITE_TIMEOUT" envDefault:"30s"`
	}
)

// New initializes and returns the configuration from environment variables and .env file
func New() (*Config, error) {
	var cfg Config

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Parse environment variables into the Config structure
	err = env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}

	return &cfg, nil
}
