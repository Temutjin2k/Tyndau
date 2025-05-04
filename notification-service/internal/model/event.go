// internal/model/event.go
package model

import "time"

// EventType представляет тип события
type EventType string

const (
	EventTypeUserRegistered EventType = "user.registered"
	EventTypeAlbumReleased  EventType = "music.album_released"
    EventTypeAlbumReleasedMass  EventType = "music.album_released_mass"
)

// Event представляет событие из NATS
type Event struct {
	Type      EventType   `json:"event_type"`
	UserID    string      `json:"user_id"`
	Email     string      `json:"email"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// UserRegisteredData содержит данные для события регистрации пользователя
type UserRegisteredData struct {
	Name string `json:"name"`
}

// AlbumReleasedData содержит данные для события релиза альбома
type AlbumReleasedData struct {
	AlbumName  string   `json:"album_name"`
	ArtistName string   `json:"artist_name"`
	Emails     []string `json:"emails,omitempty"` // Массив email-адресов для рассылки
}