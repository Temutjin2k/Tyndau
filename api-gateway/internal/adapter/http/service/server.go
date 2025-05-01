package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Temutjin2k/Tyndau/api-gateway/config"

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

	err := api.setupRoutes(ctx, grpcMux)
	if err != nil {
		return nil, err
	}

	// Default ServeMux
	mainMux := http.NewServeMux()

	// Healthcheck endpoint
	mainMux.HandleFunc("/healthcheck", api.HealthCheck)

	// gRPC-Gateway mux
	mainMux.Handle("/", grpcMux)

	api.server = &http.Server{
		Addr:           api.addr,
		Handler:        mainMux,
		ReadTimeout:    cfg.Server.HTTPServer.ReadTimeout,
		WriteTimeout:   cfg.Server.HTTPServer.WriteTimeout,
		IdleTimeout:    cfg.Server.HTTPServer.IdleTimeout,
		MaxHeaderBytes: cfg.Server.HTTPServer.MaxHeaderBytes,
	}

	return api, nil
}

func (a *API) setupRoutes(ctx context.Context, mux *runtime.ServeMux) error {
	err := a.RegisterUsergRPCHandler(ctx, mux, a.cfg.Server.UserGRPCServers.Addr)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to register user gRPC server")
		return err
	}

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
