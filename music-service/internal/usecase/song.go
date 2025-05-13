package usecase

import (
	"context"

	"github.com/Temutjin2k/Tyndau/music-service/internal/model"
)

type Song struct {
	songRepo SongRepository
}

func NewSongService(songRepo SongRepository) *Song {
	return &Song{songRepo: songRepo}
}

// Upload saves song metadata after file is uploaded to MinIO
func (s *Song) Upload(ctx context.Context, req model.Song) (model.Song, error) {
	panic("unimplemented")
}

// UploadURL generates a presigned PUT URL and returns file URL
func (s *Song) UploadURL(ctx context.Context, filename string) (string, error) {
	panic("unimplemented")
}

// GetSong fetches a song by its ID
func (s *Song) GetSong(ctx context.Context, id int64) (model.Song, error) {
	panic("unimplemented")
}

// Search searches songs by query
func (s *Song) Search(ctx context.Context, req model.SongSearch) ([]model.Song, error) {
	panic("unimplemented")
}

// Delete deletes a song by ID
func (s *Song) Delete(ctx context.Context, id int64) error {
	panic("unimplemented")
}
