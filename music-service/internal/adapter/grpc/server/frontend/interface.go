package frontend

import (
	"context"

	"github.com/Temutjin2k/Tyndau/music-service/internal/model"
)

type SongUseCase interface {
	// Upload saves song metadata after file is uploaded to MinIO
	Upload(ctx context.Context, req model.Song) (model.Song, model.UploadLink, error)

	// GetSong fetches a song by its ID
	GetSong(ctx context.Context, id int64) (model.Song, error)

	List(ctx context.Context, req model.ListRequest) ([]model.Song, error)

	Update(ctx context.Context, req model.SongUpdate) (model.Song, error)
	// Delete deletes a song by ID
	Delete(ctx context.Context, id int64) error
}
