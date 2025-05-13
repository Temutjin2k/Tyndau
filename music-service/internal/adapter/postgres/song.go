package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type SongRepository struct {
	db *pgxpool.Pool
}

func NewSongRepository(db *pgxpool.Pool) *SongRepository {
	return &SongRepository{db: db}
}
