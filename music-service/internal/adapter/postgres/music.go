package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type MusicRepository struct {
	db *pgxpool.Pool
}

func NewMusic() *MusicRepository {
	return &MusicRepository{}
}
