package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Temutjin2k/Tyndau/auth_service/config"
	grpcserver "github.com/Temutjin2k/Tyndau/auth_service/internal/adapter/grpc/server"
	"github.com/Temutjin2k/Tyndau/auth_service/internal/adapter/nats"
	"github.com/Temutjin2k/Tyndau/auth_service/internal/usecase"
	"github.com/Temutjin2k/Tyndau/auth_service/pkg/grpcconn"

	"github.com/rs/zerolog"
)

const serviceName = "user-service"

type App struct {
	grpcServer   *grpcserver.API
	natsProducer *nats.Producer

	logger *zerolog.Logger
}

func New(ctx context.Context, cfg *config.Config, logger *zerolog.Logger) (*App, error) {
	logger.Info().Str("service", serviceName).Msg("starting service")

	userConn, err := grpcconn.New(cfg.GRPCServices.UserGRPCService.Addr)
	if err != nil {
		logger.Error().Err(err).Str("service to connect", "user_service").Msg("failed to create grpc connention")
		return nil, err
	}

	// User microservice
	userProvider := usecase.NewUserProvider(userConn)

	// Nats natsProducer
	natsProducer, err := nats.NewProducer(cfg.NatsProducerConfig)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect Nats")
		return nil, err
	}

	// Mail Provider
	mailProvider := usecase.NewMail(natsProducer)

	// Token service (JWT)
	tokenService := usecase.NewJwtManager(cfg.JWT.Secret)

	authUseCase := usecase.NewAuth(userProvider, mailProvider, tokenService, logger)

	grpcServer := grpcserver.New(cfg.Server.GRPCServer, logger, authUseCase)

	app := &App{
		grpcServer:   grpcServer,
		natsProducer: natsProducer,

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
