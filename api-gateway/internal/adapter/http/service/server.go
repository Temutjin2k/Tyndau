package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Temutjin2k/Tyndau/api-gateway/config"
	m "github.com/Temutjin2k/Tyndau/api-gateway/internal/adapter/http/service/middleware"
	authProto "github.com/Temutjin2k/TyndauProto/gen/go/auth"
	musicProto "github.com/Temutjin2k/TyndauProto/gen/go/music"
	userProto "github.com/Temutjin2k/TyndauProto/gen/go/user"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
)

const serverIPAddress = "127.0.0.1:%d" // Changed to 0.0.0.0 for external access

type API struct {
	server *http.Server
	cfg    *config.Config
	logger *zerolog.Logger
	addr   string
}

func NewAPI(ctx context.Context, cfg *config.Config, logger *zerolog.Logger) (*API, error) {
	grpcMux := runtime.NewServeMux()

	api := &API{
		cfg:    cfg,
		addr:   fmt.Sprintf(serverIPAddress, cfg.Server.HTTPServer.Port),
		logger: logger,
	}

	err := api.setupGRPCRoutes(ctx, grpcMux)
	if err != nil {
		return nil, err
	}

	// Default ServeMux
	mux := http.NewServeMux()

	// Healthcheck
	mux.HandleFunc("/healthcheck", api.HealthCheck)

	// Static files
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// gRPC-Gateway mux
	mux.Handle("/", grpcMux)

	api.server = &http.Server{
		Addr:           api.addr,
		Handler:        m.LoggingMiddleware(m.CorsMiddleware(mux), logger),
		ReadTimeout:    cfg.Server.HTTPServer.ReadTimeout,
		WriteTimeout:   cfg.Server.HTTPServer.WriteTimeout,
		IdleTimeout:    cfg.Server.HTTPServer.IdleTimeout,
		MaxHeaderBytes: cfg.Server.HTTPServer.MaxHeaderBytes,
	}

	return api, nil
}

// Setup all routes
func (a *API) setupGRPCRoutes(ctx context.Context, mux *runtime.ServeMux) error {

	a.logger.Debug().Str("auth gRPC address", a.cfg.Server.AuthGRPCServer.Addr).Msg("Trying to register auth gRPC service")
	err := a.RegisterGRPCHandler(ctx, mux, a.cfg.Server.AuthGRPCServer.Addr, authProto.RegisterAuthHandler)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to register auth gRPC server")
		return err
	}

	a.logger.Debug().Str("user gRPC address", a.cfg.Server.UserGRPCServers.Addr).Msg("Trying to register user gRPC service")
	err = a.RegisterGRPCHandler(ctx, mux, a.cfg.Server.UserGRPCServers.Addr, userProto.RegisterUserHandler)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to register user gRPC server")
		return err
	}

	a.logger.Debug().Str("music gRPC address", a.cfg.Server.MusicGRPCServer.Addr).Msg("Trying to register music gRPC service")
	err = a.RegisterGRPCHandler(ctx, mux, a.cfg.Server.MusicGRPCServer.Addr, musicProto.RegisterMusicHandler)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to register music gRPC server")
		return err
	}

	a.logger.Debug().Msg("All grpc services registered")
	return nil
}

func (a *API) Stop() error {
	// Creating a context with timeout for graceful shutdown
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.logger.Info().Msg("HTTP server shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := a.server.Shutdown(ctx)
	if err != nil {
		a.logger.Error().Err(err).Msg("Server forced to shutdown due to error")
		return err
	}

	a.logger.Info().Msg("HTTP server stopped successfully")
	return nil
}

func (a *API) Run(errCh chan<- error) {
	go func() {
		a.logger.Info().Str("Address", a.addr).Msg("Server initialized successfully, starting HTTP server")
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal().Err(err).Msg("Error starting HTTP server")
		}
	}()
}

func (a *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	responce := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"address": a.addr,
			"version": a.cfg.Version,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responce)
}
