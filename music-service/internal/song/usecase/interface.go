package usecase

import (
	"context"
	"github.com/Temutjin2k/Tyndau/music-service/internal/song/entity"
)

type SongService interface {
	Create(ctx context.Context, in *entity.Song) (*entity.Song, error)
	Get(ctx context.Context, id int64) (*entity.Song, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Song, error)
	Update(ctx context.Context, in *entity.Song) error
	Delete(ctx context.Context, id int64) error
}

type SongRepository interface {
	Create(ctx context.Context, in *entity.Song) (*entity.Song, error)
	Get(ctx context.Context, id int64) (*entity.Song, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Song, error)
	Update(ctx context.Context, in *entity.Song) error
	Delete(ctx context.Context, id int64) error
}
