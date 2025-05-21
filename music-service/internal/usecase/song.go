package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Temutjin2k/Tyndau/music-service/internal/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Song struct {
	repo     SongRepository
	cache    Cache
	uploader UploaderV3

	logger *zerolog.Logger
}

func NewSongService(songRepo SongRepository, cache Cache, uploader UploaderV3, logger *zerolog.Logger) *Song {
	return &Song{
		repo:     songRepo,
		cache:    cache,
		uploader: uploader,

		logger: logger,
	}
}

// Upload saves song metadata after file is uploaded to MinIO
func (s *Song) Upload(ctx context.Context, req model.Song) (model.Song, model.UploadLink, error) {
	logger := s.logger.With().Str("filename", req.Filename).Logger()

	// 1. Генерируем presigned URL и публичный URL для загрузки файла
	objectKey := generateObjectKey(req.Filename)

	uploadURL, publicURL, err := s.uploader.GenerateSongURLs(ctx, objectKey)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate presigned URLs")
		return model.Song{}, model.UploadLink{}, fmt.Errorf("failed to generate presigned URLs: %w", err)
	}

	// 2. Записываем публичный URL в модель как FileURL
	req.FileURL = publicURL

	// 3. Создаем запись песни с метаданными и FileURL
	out, err := s.repo.Create(ctx, &req)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create song record")
		return model.Song{}, model.UploadLink{}, fmt.Errorf("failed to create song record: %w", err)
	}

	// 4. Сбрасываем кэш по песне
	if err := s.cache.Delete(ctx, out.ID); err != nil {
		logger.Warn().Err(err).Msg("failed to invalidate cache after upload")
	}

	logger.Info().Msg("song record created with presigned URLs")

	// 5. Возвращаем данные и ссылки для загрузки
	return *out, model.UploadLink{
		UploadURL: uploadURL,
		FileURL:   publicURL,
	}, nil
}

// GetSong fetches a song by its ID
func (s *Song) GetSong(ctx context.Context, id int64) (model.Song, error) {
	logger := s.logger.With().Int64("song_id", id).Logger()

	// Check cache first
	if hit, err := s.cache.Get(ctx, id); hit != nil {
		if err != nil {
			logger.Warn().Err(err).Msg("cache lookup error")
		} else {
			logger.Debug().Msg("song retrieved from cache")
			return *hit, nil
		}
	}

	// Get from repository
	res, err := s.repo.Get(ctx, id)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get song from repository")
		return model.Song{}, fmt.Errorf("failed to get song: %w", err)
	}
	if res == nil {
		logger.Info().Msg("song not found")
		return model.Song{}, fmt.Errorf("song with id %d not found", id)
	}

	// Update cache
	if err := s.cache.Set(ctx, res); err != nil {
		logger.Warn().Err(err).Msg("failed to cache song")
	}

	logger.Debug().Msg("song retrieved from repository")
	return *res, nil
}

// Search searches songs by query
func (s *Song) List(ctx context.Context, req model.ListRequest) ([]model.Song, error) {
	logger := s.logger.With().
		Str("query", req.Query).
		Int32("limit", req.Limit).
		Int32("offset", req.Offset).
		Logger()

	songs, err := s.repo.List(ctx, req)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list songs")
		return nil, fmt.Errorf("failed to list songs: %w", err)
	}

	logger.Debug().Int("count", len(songs)).Msg("songs listed successfully")
	return songs, nil
}

// Update updates song metadata
func (s *Song) Update(ctx context.Context, req model.SongUpdate) (model.Song, error) {
	logger := s.logger.With().Int64("song_id", req.ID).Logger()

	// 1. Get existing song
	existing, err := s.GetSong(ctx, req.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get song for update")
		return model.Song{}, fmt.Errorf("failed to get song: %w", err)
	}

	// 2. Apply partial updates
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Artist != nil {
		existing.Artist = *req.Artist
	}
	if req.Album != nil {
		existing.Album = *req.Album
	}
	if req.Genre != nil {
		existing.Genre = *req.Genre
	}
	if req.DurationSeconds != nil {
		existing.DurationSeconds = *req.DurationSeconds
	}
	if req.ReleaseDate != nil {
		// Parse the string date into time.Time
		parsedDate, err := time.Parse(time.RFC3339, *req.ReleaseDate)
		if err != nil {
			return model.Song{}, fmt.Errorf("invalid release date format: %w", err)
		}
		existing.ReleaseDate = parsedDate
	}

	// 3. Save to repository
	if err := s.repo.Update(ctx, &existing); err != nil {
		logger.Error().Err(err).Msg("failed to update song in repository")
		return model.Song{}, fmt.Errorf("failed to update song: %w", err)
	}

	// 4. Invalidate cache
	if err := s.cache.Delete(ctx, req.ID); err != nil {
		logger.Warn().Err(err).Msg("failed to invalidate cache after update")
	}

	logger.Info().Msg("song updated successfully")
	return existing, nil
}

// Delete deletes a song by ID
func (s *Song) Delete(ctx context.Context, id int64) error {
	logger := s.logger.With().Int64("song_id", id).Logger()

	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error().Err(err).Msg("failed to delete song from repository")
		return fmt.Errorf("failed to delete song: %w", err)
	}

	if err := s.cache.Delete(ctx, id); err != nil {
		logger.Warn().Err(err).Msg("failed to invalidate cache after deletion")
	}

	logger.Info().Msg("song deleted successfully")
	return nil
}

func generateObjectKey(fileName string) string {
	id := uuid.New()
	return fmt.Sprintf("songs/%s_%s", id.String(), fileName)
}
