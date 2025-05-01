package server

import "github.com/Temutjin2k/Tyndau/user_service/internal/adapter/grpc/server/frontend"

type UserUseCase interface {
	frontend.UserUseCase
}
