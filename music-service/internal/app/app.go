package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Temutjin2k/Tyndau/music-service/config"
	grpcserver "github.com/Temutjin2k/Tyndau/music-service/internal/adapter/grpc/server"
	postgresrepo "github.com/Temutjin2k/Tyndau/music-service/internal/adapter/postgres"
	"github.com/Temutjin2k/Tyndau/music-service/internal/usecase"
	"github.com/Temutjin2k/Tyndau/music-service/pkg/postgres"
	"github.com/rs/zerolog"
)

const serviceName = "music-service"

type App struct {
	grpcServer *grpcserver.API
	postgresDB *postgres.PostgreDB

	logger *zerolog.Logger
}

func New(ctx context.Context, cfg *config.Config, logger *zerolog.Logger) (*App, error) {

	postgresDB, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to Postgres")
		return nil, fmt.Errorf("postgres: %w", err)
	}

	_ = postgresrepo.NewSongRepository(postgresDB.Pool)
	songUseCase := usecase.NewSongService(nil)
	grpcServer := grpcserver.New(cfg.Server.GRPCServer, logger, songUseCase)

	app := &App{
		grpcServer: grpcServer,
		postgresDB: postgresDB,

		logger: logger,
	}

	return app, nil
}

func (a *App) Close(ctx context.Context) {
	err := a.grpcServer.Stop(ctx)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to stop http server")
	}
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
