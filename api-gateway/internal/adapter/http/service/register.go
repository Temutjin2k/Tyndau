package service

import (
	"context"
	"time"

	"github.com/Temutjin2k/Tyndau/api-gateway/pkg/grpcconn"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func (a *API) RegisterGRPCHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	endpoint string,
	register func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error,
) error {
	timeOutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	conn, err := grpcconn.New(a.cfg.Server.UserGRPCServers.Addr)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to create connection")
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", a.cfg.Server.UserGRPCServers.Addr, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", a.cfg.Server.UserGRPCServers.Addr, cerr)
			}
		}()
	}()

	if err := register(timeOutCtx, mux, conn); err != nil {
		a.logger.Error().Err(err).Msg("Failed to register grpc handler")
		return err
	}

	return nil
}
