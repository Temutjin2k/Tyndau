package service

import (
	"context"
	"time"

	"github.com/Temutjin2k/Tyndau/api-gateway/pkg/grpcconn"
	userProto "github.com/Temutjin2k/TyndauProto/gen/go/user"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/grpclog"
)

// RegisterUsergRPCHandler registers User Handler, also opens and closes the connection
func (a *API) RegisterUsergRPCHandler(ctx context.Context, mux *runtime.ServeMux, endpoint string) error {
	timeOutCtx, userCancel := context.WithTimeout(ctx, 3*time.Second)
	defer userCancel()

	userServiceConn, err := grpcconn.New(a.cfg.Server.UserGRPCServers.Addr)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to create connection")
		return err
	}
	defer func() {
		if err != nil {
			if cerr := userServiceConn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", a.cfg.Server.UserGRPCServers.Addr, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := userServiceConn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", a.cfg.Server.UserGRPCServers.Addr, cerr)
			}
		}()
	}()

	if err := userProto.RegisterUserHandler(timeOutCtx, mux, userServiceConn); err != nil {
		a.logger.Error().Err(err).Msg("Failed to register User handler")
		return err
	}

	return nil
}
