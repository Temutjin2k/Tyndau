package usecase

import (
	"context"

	"github.com/Temutjin2k/Tyndau/music-service/internal/model"
)

type Cache interface {
	Get(ctx context.Context, id int64) (*model.Song, error)
	Set(ctx context.Context, song *model.Song) error
	Delete(ctx context.Context, id int64) error
}

type SongRepository interface {
	Create(ctx context.Context, req *model.Song) (*model.Song, error)
	Get(ctx context.Context, id int64) (*model.Song, error)
	List(ctx context.Context, req model.ListRequest) ([]model.Song, error)
	Delete(ctx context.Context, id int64) error
}
