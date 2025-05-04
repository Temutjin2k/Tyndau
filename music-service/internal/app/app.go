package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Temutjin2k/Tyndau/music-service/config"
	grpcserver "github.com/Temutjin2k/Tyndau/music-service/internal/adapter/grpc/server"
	"github.com/Temutjin2k/Tyndau/music-service/internal/adapter/nats"
	"github.com/rs/zerolog"
)

const serviceName = "music-service"

type App struct {
	grpcServer   *grpcserver.API
	natsProducer *nats.Producer

	logger *zerolog.Logger
}

func New(ctx context.Context, cfg *config.Config, logger *zerolog.Logger) (*App, error) {
	app := &App{
		logger: logger,
	}

	return app, nil
}

func (a *App) Close(ctx context.Context) {
	err := a.grpcServer.Stop(ctx)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to stop http server")
	}

	a.natsProducer.Close()

}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()

	a.grpcServer.Run(ctx, errCh)

	a.logger.Info().Str("name", serviceName).Msg("service started")

	// Waiting signal
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun

	case s := <-shutdownCh:
		a.logger.Info().Str("received signal", s.String()).Msg("Shutting down application")

		a.Close(ctx)
		a.logger.Info().Msg("Application stopped!")
	}

	return nil
}
