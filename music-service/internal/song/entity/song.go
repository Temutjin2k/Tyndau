package entity

import "time"

type Song struct {
	ID          int64     `json:"id"            db:"id"`
	Title       string    `json:"title"         db:"title"`
	Artist      string    `json:"artist"        db:"artist"`
	Album       string    `json:"album,omitempty" db:"album"`
	Genre       string    `json:"genre,omitempty" db:"genre"`
	DurationSec int32     `json:"duration_sec"  db:"duration_sec"`
	FileURL     string    `json:"file_url,omitempty" db:"file_url"`
	ReleasedAt  time.Time `json:"released_at"   db:"released_at"`
	CreatedAt   time.Time `json:"created_at"    db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"    db:"updated_at"`
}
