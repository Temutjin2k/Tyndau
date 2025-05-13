package model

import "time"

type Song struct {
	ID              int64
	Title           string
	Artist          string
	Album           string
	Genre           string
	ReleaseDate     time.Time
	DurationSeconds int32
	FileURL         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type SongSearch struct {
	Query  string
	Limit  int32
	Offset int32
}
