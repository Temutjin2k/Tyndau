package usecase

import (
	"context"
	"fmt"

	"github.com/Temutjin2k/Tyndau/music-service/internal/model"
	"github.com/rs/zerolog"
)

type Song struct {
	repo   SongRepository
	cache  Cache
	logger zerolog.Logger
}

func NewSongService(songRepo SongRepository, cache Cache, logger *zerolog.Logger) *Song {
	return &Song{
		repo:   songRepo,
		cache:  cache,
		logger: logger.With().Str("component", "song_usecase").Logger(),
	}
}

// Upload saves song metadata after file is uploaded to MinIO
func (s *Song) Upload(ctx context.Context, req model.Song) (model.Song, error) {
	logger := s.logger.With().Int64("song_id", req.ID).Logger()

	out, err := s.repo.Create(ctx, &req)
	if err != nil {
		logger.Error().Err(err).Msg("failed to upload song")
		return model.Song{}, fmt.Errorf("failed to upload song: %w", err)
	}

	if err := s.cache.Delete(ctx, out.ID); err != nil {
		logger.Warn().Err(err).Msg("failed to invalidate cache after upload")
	}

	logger.Info().Msg("song uploaded successfully")
	return *out, nil
}

// UploadURL generates a presigned PUT URL and returns file URL
func (s *Song) UploadURL(ctx context.Context, filename string) (string, error) {
	s.logger.Warn().Msg("UploadURL not implemented")
	return "", fmt.Errorf("UploadURL not implemented")
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
