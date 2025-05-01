package server

import (
	"context"
	"fmt"
	"net"

	"github.com/Temutjin2k/Tyndau/user_service/config"
	"github.com/Temutjin2k/Tyndau/user_service/internal/adapter/grpc/server/frontend"
	userpb "github.com/Temutjin2k/TyndauProto/gen/go/user"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type API struct {
	s      *grpc.Server
	cfg    config.GRPCServer
	addr   string
	logger *zerolog.Logger

	userUsecase UserUseCase
}

func New(cfg config.GRPCServer, logger *zerolog.Logger, userUsecase UserUseCase) *API {
	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)

	if logger != nil {
		logger.Info().
			Str("address", addr).
			Msg("gRPC server instance created")
	}

	return &API{
		cfg:    cfg,
		addr:   addr,
		logger: logger,

		userUsecase: userUsecase,
	}
}

func (a *API) Run(ctx context.Context, errCh chan<- error) {
	go func() {
		a.logger.Info().
			Str("address", a.addr).
			Msg("gRPC server starting")

		if err := a.run(ctx); err != nil {
			a.logger.Error().
				Err(err).
				Msg("gRPC server failed to start")
			errCh <- fmt.Errorf("can't start grpc server: %w", err)
		}
	}()
}

func (a *API) Stop(ctx context.Context) error {
	if a.s == nil {
		a.logger.Debug().Msg("gRPC server already stopped (nil server instance)")
		return nil
	}

	stopped := make(chan struct{})
	go func() {
		a.logger.Info().Msg("gRPC server shutting down gracefully")
		a.s.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		a.logger.Warn().Msg("gRPC server force stopped due to context cancellation")
		a.s.Stop()
	case <-stopped:
		a.logger.Info().Msg("gRPC server stopped gracefully")
	}

	return nil
}

func (a *API) run(ctx context.Context) error {
	a.s = grpc.NewServer(a.setOptions(ctx)...)

	// Register services
	userServer := frontend.NewUser(a.userUsecase, a.logger)

	userpb.RegisterUserServer(a.s, userServer)
	reflection.Register(a.s)

	a.logger.Debug().Msg("gRPC services registered")

	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		a.logger.Error().
			Err(err).
			Str("address", a.addr).
			Msg("failed to create listener")
		return fmt.Errorf("failed to create listener: %w", err)
	}

	a.logger.Info().
		Str("address", a.addr).
		Msg("gRPC server started listening")

	if err := a.s.Serve(listener); err != nil {
		a.logger.Error().
			Err(err).
			Msg("gRPC server failed to serve")
		return fmt.Errorf("failed to serve grpc: %w", err)
	}

	return nil
}
