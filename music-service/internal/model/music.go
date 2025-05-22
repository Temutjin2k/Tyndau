package model

import "time"

type Song struct {
	ID              int64
	Title           string
	Artist          string
	Album           string
	Genre           string
	Filename        string
	FileURL         string
	DurationSeconds int32
	ReleaseDate     time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// SongUpdate represents fields that can be updated in a song
type SongUpdate struct {
	ID              int64     // Required - song ID to update
	Title           *string   `json:"title,omitempty"`            // Optional new title
	Artist          *string   `json:"artist,omitempty"`           // Optional new artist
	Album           *string   `json:"album,omitempty"`            // Optional new album
	Genre           *string   `json:"genre,omitempty"`            // Optional new genre
	DurationSeconds *int32    `json:"duration_seconds,omitempty"` // Optional new duration
	ReleaseDate     *string   `json:"release_date,omitempty"`     // Optional new release date
	FileURL         *string   `json:"file_url,omitempty"`         // Optional new file URL
	UpdatedAt       time.Time // Will be set automatically
}

type ListRequest struct {
	Query  string
	Limit  int32
	Offset int32
}

type UploadLink struct {
	UploadURL string
	FileURL   string
}
