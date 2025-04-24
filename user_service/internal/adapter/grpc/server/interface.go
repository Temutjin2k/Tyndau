package server

import "user_service/internal/adapter/grpc/server/frontend"

type UserUseCase interface {
	frontend.UserUseCase
}
