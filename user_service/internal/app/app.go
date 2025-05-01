package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Temutjin2k/Tyndau/user_service/config"
	grpcserver "github.com/Temutjin2k/Tyndau/user_service/internal/adapter/grpc/server"
	postgresrepo "github.com/Temutjin2k/Tyndau/user_service/internal/adapter/postgres"
	"github.com/Temutjin2k/Tyndau/user_service/internal/usecase"
	"github.com/Temutjin2k/Tyndau/user_service/pkg/postgres"

	"github.com/rs/zerolog"
)

const serviceName = "user-service"

type App struct {
	grpcServer *grpcserver.API
	postgresDB *postgres.PostgreDB

	logger *zerolog.Logger
}

func New(ctx context.Context, cfg *config.Config, logger *zerolog.Logger) (*App, error) {
	logger.Info().Str("service", serviceName).Msg("starting service")

	logger.Info().Msg("connecting to Postgres")

	postgresDB, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to Postgres")
		return nil, fmt.Errorf("postgres: %w", err)
	}

	logger.Info().Msg("Postgres connection established")

	userRepo := postgresrepo.NewUserRepository(postgresDB.Pool)
	userUseCase := usecase.NewUser(userRepo, logger)
	grpcServer := grpcserver.New(cfg.Server.GRPCServer, logger, userUseCase)

	app := &App{
		grpcServer: grpcServer,
		postgresDB: postgresDB,
		logger:     logger,
	}

	return app, nil
}

func (a *App) Close(ctx context.Context) {
	err := a.grpcServer.Stop(ctx)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to stop http server")
	}

	// Closing postgres connection
	a.postgresDB.Pool.Close()
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()

	a.grpcServer.Run(ctx, errCh)

	a.logger.Info().Str("name", serviceName).Msg("service is running")

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
