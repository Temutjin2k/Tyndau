package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Temutjin2k/Tyndau/music-service/internal/song/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SongRepo struct{ db *pgxpool.Pool }

func NewSongRepo(db *pgxpool.Pool) *SongRepo { return &SongRepo{db: db} }

// --- CRUD ---

func (r *SongRepo) Create(ctx context.Context, in *entity.Song) (*entity.Song, error) {
	const q = `INSERT INTO songs (title,artist,album,genre,duration_sec,file_url,released_at,created_at,updated_at)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,now(),now())
	           RETURNING id,created_at,updated_at`
	err := r.db.QueryRow(ctx, q,
		in.Title, in.Artist, in.Album, in.Genre,
		in.DurationSec, in.FileURL, in.ReleasedAt).
		Scan(&in.ID, &in.CreatedAt, &in.UpdatedAt)
	return in, err
}

func (r *SongRepo) Get(ctx context.Context, id int64) (*entity.Song, error) {
	const q = `SELECT id,title,artist,album,genre,duration_sec,file_url,
	                released_at,created_at,updated_at
	           FROM songs WHERE id=$1`
	var s entity.Song
	err := r.db.QueryRow(ctx, q, id).Scan(
		&s.ID, &s.Title, &s.Artist, &s.Album, &s.Genre, &s.DurationSec,
		&s.FileURL, &s.ReleasedAt, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &s, err
}

func (r *SongRepo) List(ctx context.Context, limit, offset int) ([]*entity.Song, error) {
	const q = `SELECT id,title,artist,album,genre,duration_sec,file_url,
	                released_at,created_at,updated_at
	           FROM songs ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*entity.Song
	for rows.Next() {
		var s entity.Song
		if err = rows.Scan(&s.ID, &s.Title, &s.Artist, &s.Album, &s.Genre,
			&s.DurationSec, &s.FileURL, &s.ReleasedAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, &s)
	}
	return res, rows.Err()
}

func (r *SongRepo) Update(ctx context.Context, in *entity.Song) error {
	in.UpdatedAt = time.Now()
	const q = `UPDATE songs
	           SET title=$2,artist=$3,album=$4,genre=$5,duration_sec=$6,file_url=$7,
	               released_at=$8,updated_at=$9
	           WHERE id=$1`
	_, err := r.db.Exec(ctx, q, in.ID, in.Title, in.Artist, in.Album, in.Genre,
		in.DurationSec, in.FileURL, in.ReleasedAt, in.UpdatedAt)
	return err
}

func (r *SongRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM songs WHERE id=$1`, id)
	return err
}
