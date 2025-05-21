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

type ListRequest struct {
	Query  string
	Limit  int32
	Offset int32
}

type UploadLink struct {
	UploadURL string
	FileURL   string
}
