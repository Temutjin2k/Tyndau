package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Temutjin2k/Tyndau/api-gateway/config"
	httpserver "github.com/Temutjin2k/Tyndau/api-gateway/internal/adapter/http/service"
	"github.com/rs/zerolog"
)

const serviceName string = "api-gateway"

type App struct {
	httpServer *httpserver.API
	logger     *zerolog.Logger
}

func New(ctx context.Context, cfg *config.Config, logger *zerolog.Logger) (*App, error) {
	httpServer, err := httpserver.NewAPI(ctx, cfg, logger)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create http server")
		return nil, err
	}

	app := &App{
		httpServer: httpServer,
		logger:     logger,
	}
	return app, nil
}

func (a *App) Close(ctx context.Context) {
	err := a.httpServer.Stop()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to stop http server")
	}
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()

	// Running http server
	a.httpServer.Run(errCh)

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
