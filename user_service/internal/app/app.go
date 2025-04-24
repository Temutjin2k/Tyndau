package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"user_service/config"
	grpcserver "user_service/internal/adapter/grpc/server"
	postgresrepo "user_service/internal/adapter/postgres"
	"user_service/internal/usecase"
	"user_service/pkg/postgres"
)

const serviceName = "user-service"

type App struct {
	grpcServer *grpcserver.API
	postgresDB *postgres.PostgreDB
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Printf("starting %v service\n", serviceName)

	log.Println("connecting to postgres")

	postgresDB, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}
	log.Println("connection established")

	userRepo := postgresrepo.NewUserRepository(postgresDB.Pool)

	userUseCase := usecase.NewUser(userRepo)

	grpcServer := grpcserver.New(cfg.Server.GRPCServer, userUseCase)

	app := &App{
		grpcServer: grpcServer,
		postgresDB: postgresDB,
	}

	return app, nil
}

func (a *App) Close(ctx context.Context) {
	err := a.grpcServer.Stop(ctx)
	if err != nil {
		log.Println("failed to shutdown gRPC service", err)
	}

	// Closing postgres connection
	a.postgresDB.Pool.Close()
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()

	a.grpcServer.Run(ctx, errCh)

	log.Printf("service %v started\n", serviceName)

	// Waiting signal
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun

	case s := <-shutdownCh:
		log.Printf("received signal: %v. Running graceful shutdown...\n", s)

		a.Close(ctx)
		log.Println("graceful shutdown completed!")
	}

	return nil
}
