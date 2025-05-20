package frontend

import (
	"context"

	"github.com/Temutjin2k/Tyndau/music-service/internal/model"
)

type SongUseCase interface {
	// Upload saves song metadata after file is uploaded to MinIO
	Upload(ctx context.Context, req model.Song) (model.Song, error)

	// UploadURL generates a presigned PUT URL and returns file URL
	UploadURL(ctx context.Context, filename string) (string, error)

	// GetSong fetches a song by its ID
	GetSong(ctx context.Context, id int64) (model.Song, error)

	List(ctx context.Context, req model.ListRequest) ([]model.Song, error)

	// Delete deletes a song by ID
	Delete(ctx context.Context, id int64) error
}
