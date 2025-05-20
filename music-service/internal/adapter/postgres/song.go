package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Temutjin2k/Tyndau/music-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SongRepository struct {
	db *pgxpool.Pool
}

func NewSongRepository(db *pgxpool.Pool) *SongRepository {
	return &SongRepository{db: db}
}

func (r *SongRepository) Create(ctx context.Context, in *model.Song) (*model.Song, error) {
	const q = `
		INSERT INTO songs (
			title, artist, album, genre, duration_sec, file_url, released_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, now(), now()
		)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx, q,
		in.Title,
		in.Artist,
		in.Album,
		in.Genre,
		in.DurationSeconds,
		in.FileURL,
		in.ReleaseDate,
	).Scan(
		&in.ID,
		&in.CreatedAt,
		&in.UpdatedAt,
	)

	return in, err
}

func (r *SongRepository) Get(ctx context.Context, id int64) (*model.Song, error) {
	const q = `
		SELECT 
			id, title, artist, album, genre, duration_sec, file_url,
			released_at, created_at, updated_at
		FROM 
			songs 
		WHERE 
			id = $1
	`

	var s model.Song
	err := r.db.QueryRow(ctx, q, id).Scan(
		&s.ID,
		&s.Title,
		&s.Artist,
		&s.Album,
		&s.Genre,
		&s.DurationSeconds,
		&s.FileURL,
		&s.ReleaseDate,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &s, err
}

func (r *SongRepository) List(ctx context.Context, req model.ListRequest) ([]model.Song, error) {
	const q = `
		SELECT 
			id, title, artist, album, genre, duration_sec, file_url,
			released_at, created_at, updated_at
		FROM 
			songs 
		ORDER BY 
			id 
		LIMIT 
			$1 
		OFFSET 
			$2
	`

	rows, err := r.db.Query(ctx, q, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Song
	for rows.Next() {
		var s model.Song

		if err = rows.Scan(
			&s.ID,
			&s.Title,
			&s.Artist,
			&s.Album,
			&s.Genre,
			&s.DurationSeconds,
			&s.FileURL,
			&s.ReleaseDate,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}

		res = append(res, s)
	}

	return res, rows.Err()
}

func (r *SongRepository) Update(ctx context.Context, in *model.Song) error {
	in.UpdatedAt = time.Now()

	const q = `
		UPDATE 
			songs
		SET 
			title = $2,
			artist = $3,
			album = $4,
			genre = $5,
			duration_sec = $6,
			file_url = $7,
			released_at = $8,
			updated_at = $9
		WHERE 
			id = $1
	`

	_, err := r.db.Exec(
		ctx,
		q,
		in.ID,
		in.Title,
		in.Artist,
		in.Album,
		in.Genre,
		in.DurationSeconds,
		in.FileURL,
		in.ReleaseDate,
		in.UpdatedAt,
	)

	return err
}

func (r *SongRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM songs WHERE id = $1`, id)
	return err
}
