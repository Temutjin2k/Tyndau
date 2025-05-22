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

	ctx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

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
	// TODO search, sorting, pagination
	const q = `
		SELECT 
			id, title, artist, album, genre, duration_sec, file_url,
			released_at, created_at, updated_at
		FROM 
			songs 
		ORDER BY 
			id 
	`

	rows, err := r.db.Query(ctx, q)
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

// Update updates song metadata in database
func (r *SongRepository) Update(ctx context.Context, in *model.Song) error {
	in.UpdatedAt = time.Now()

	const q = `
        UPDATE songs
        SET 
            title = COALESCE($2, title),
            artist = COALESCE($3, artist),
            album = COALESCE($4, album),
            genre = COALESCE($5, genre),
            duration_sec = COALESCE($6, duration_sec),
            released_at = COALESCE($7, released_at),
            updated_at = $8
        WHERE id = $1
    `

	_, err := r.db.Exec(
		ctx, q,
		in.ID,
		in.Title,
		in.Artist,
		in.Album,
		in.Genre,
		in.DurationSeconds,
		in.ReleaseDate,
		in.UpdatedAt,
	)

	return err
}

func (r *SongRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM songs WHERE id = $1`, id)
	return err
}
