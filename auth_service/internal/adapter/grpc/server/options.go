package server

import (
	"context"

	i "github.com/Temutjin2k/Tyndau/auth_service/internal/adapter/grpc/server/interseptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func (a *API) setOptions(ctx context.Context) []grpc.ServerOption {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(i.UnaryLoggingInterceptor(a.logger)),

		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge:      a.cfg.MaxConnectionAge,
			MaxConnectionAgeGrace: a.cfg.MaxConnectionAgeGrace,
		}),
		grpc.MaxRecvMsgSize(a.cfg.MaxRecvMsgSizeMiB * (1024 * 1024)), // MaxRecvSize * 1 MB
	}

	return opts
}
