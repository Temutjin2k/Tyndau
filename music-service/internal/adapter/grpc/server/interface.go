package server

import "github.com/Temutjin2k/Tyndau/music-service/internal/adapter/grpc/server/frontend"

type MusicUseCase interface {
	frontend.SongUseCase
}
