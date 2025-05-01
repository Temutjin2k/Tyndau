package server

import "github.com/Temutjin2k/Tyndau/auth_service/internal/adapter/grpc/server/frontend"

type AuthUseCase interface {
	frontend.AuthUseCase
}
