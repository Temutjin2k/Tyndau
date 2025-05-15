package model

type Author struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

type Album struct {
	ID          int64  `json:"id"`
	AuthorID    int64  `json:"author_id"`
	Name        string `json:"name"`
	Year        int    `json:"year"`
	Description string `json:"description"`
}

type Genre struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TrackStatus string

const (
	StatusDraft     TrackStatus = "draft"
	StatusPublished TrackStatus = "published"
)

type SongSearch struct {
	Query  string
	Limit  int32
	Offset int32
}
