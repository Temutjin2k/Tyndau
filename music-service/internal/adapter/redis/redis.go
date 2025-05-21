package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Temutjin2k/Tyndau/music-service/config"
	"github.com/Temutjin2k/Tyndau/music-service/internal/model"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client  *redis.Client
	ttl     time.Duration
	keyPref string
}

// New создает новый экземпляр Redis кэша с проверкой подключения
func New(ctx context.Context, cfg config.Redis) (*redisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Checking connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &redisCache{
		client:  rdb,
		ttl:     cfg.TTL,
		keyPref: "song:",
	}, nil
}

// Ping check server availability
func (r *redisCache) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	return err
}

func (r *redisCache) key(id int64) string {
	return fmt.Sprintf("%s%d", r.keyPref, id)
}

func (r *redisCache) Get(ctx context.Context, id int64) (*model.Song, error) {
	key := r.key(id)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var sr SongRedis
	if err := json.Unmarshal([]byte(data), &sr); err != nil {
		return nil, err
	}

	return sr.ToSong(), nil
}

func (r *redisCache) Set(ctx context.Context, s *model.Song) error {
	sr := FromSong(s)

	data, err := json.Marshal(sr)
	if err != nil {
		return err
	}

	key := r.key(s.ID)
	log.Println(">>> SET REDIS:", key)

	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		log.Println(">>> REDIS SET ERROR:", err)
		return err
	}
	return nil
}

func (r *redisCache) Delete(ctx context.Context, id int64) error {
	return r.client.Del(ctx, r.key(id)).Err()
}
