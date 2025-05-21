package server

import (
	"context"

	"github.com/Temutjin2k/Tyndau/music-service/internal/adapter/grpc/server/interseptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func (a *API) setOptions(ctx context.Context) []grpc.ServerOption {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(interseptor.UnaryLoggingInterceptor(a.logger)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge:      a.cfg.MaxConnectionAge,
			MaxConnectionAgeGrace: a.cfg.MaxConnectionAgeGrace,
		}),
		grpc.MaxRecvMsgSize(a.cfg.MaxRecvMsgSizeMiB * (1024 * 1024)), // MaxRecvSize * 1 MB
	}

	return opts
}
