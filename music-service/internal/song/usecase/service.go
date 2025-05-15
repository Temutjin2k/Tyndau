package usecase

import (
	"context"
	"github.com/Temutjin2k/Tyndau/music-service/internal/song/cache/redis"
	"github.com/Temutjin2k/Tyndau/music-service/internal/song/entity"
	"log"
)

type Song struct {
	repo  SongRepository
	cache redis.Cache
}

func NewSongService(r SongRepository, c redis.Cache) SongService {
	return &Song{repo: r, cache: c}
}

func (s *Song) Create(ctx context.Context, in *entity.Song) (*entity.Song, error) {
	out, err := s.repo.Create(ctx, in)
	if err == nil {
		_ = s.cache.Delete(ctx, out.ID) // на случай повторного GET сразу после Create
	}
	return out, err
}

func (s *Song) Get(ctx context.Context, id int64) (*entity.Song, error) {
	if hit, _ := s.cache.Get(ctx, id); hit != nil {
		log.Println(">>> FROM REDIS")
		return hit, nil
	}

	res, err := s.repo.Get(ctx, id)
	if err == nil && res != nil {
		_ = s.cache.Set(ctx, res)
	}
	return res, err
}

func (s *Song) List(ctx context.Context, l, o int) ([]*entity.Song, error) {
	return s.repo.List(ctx, l, o) // кэш не нужен
}

func (s *Song) Update(ctx context.Context, in *entity.Song) error {
	if err := s.repo.Update(ctx, in); err != nil {
		return err
	}
	return s.cache.Delete(ctx, in.ID)
}

func (s *Song) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return s.cache.Delete(ctx, id)
}
