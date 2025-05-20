package usecase

import (
	"context"
	"time"

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

// Uploader defines the behavior for generating secure upload links.
type Uploader interface {
	// PresignedPutURL generates a temporary URL for uploading a file directly to storage.
	// bucket: the name of the bucket (e.g., "music")
	// objectName: the name of the file (e.g., "song123.mp3")
	// expires: how long the URL is valid (e.g., 15 minutes)
	PresignedPutURL(ctx context.Context, bucket string, objectName string, expires time.Duration) (string, error)
}

type UploaderV2 interface {
	GenerateSongURLs(ctx context.Context, songID int64, fileName string) (string, string, error)
}
