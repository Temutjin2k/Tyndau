package redis

import (
	"time"

	"github.com/Temutjin2k/Tyndau/music-service/internal/model"
)

type SongRedis struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Artist          string    `json:"artist"`
	Album           string    `json:"album"`
	Genre           string    `json:"genre"`
	FileURL         string    `json:"file_url"`
	DurationSeconds int32     `json:"duration_seconds"`
	ReleaseDate     time.Time `json:"release_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ToSong преобразует SongResis в Song
func (sr *SongRedis) ToSong() *model.Song {
	return &model.Song{
		ID:              sr.ID,
		Title:           sr.Title,
		Artist:          sr.Artist,
		Album:           sr.Album,
		Genre:           sr.Genre,
		FileURL:         sr.FileURL,
		DurationSeconds: sr.DurationSeconds,
		ReleaseDate:     sr.ReleaseDate,
		CreatedAt:       sr.CreatedAt,
		UpdatedAt:       sr.UpdatedAt,
	}
}

// FromSong преобразует Song в SongResis
func FromSong(s *model.Song) *SongRedis {
	return &SongRedis{
		ID:              s.ID,
		Title:           s.Title,
		Artist:          s.Artist,
		Album:           s.Album,
		Genre:           s.Genre,
		FileURL:         s.FileURL,
		DurationSeconds: s.DurationSeconds,
		ReleaseDate:     s.ReleaseDate,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}
