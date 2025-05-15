package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Temutjin2k/Tyndau/music-service/internal/song/entity"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, id int64) (*entity.Song, error)
	Set(ctx context.Context, s *entity.Song) error
	Delete(ctx context.Context, id int64) error
}

type redisCache struct {
	client  *redis.Client
	ttl     time.Duration
	keyPref string
}

func New(addr, password string, db int, ttl time.Duration) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &redisCache{
		client:  rdb,
		ttl:     ttl,
		keyPref: "song:",
	}
}

func (r *redisCache) key(id int64) string {
	return fmt.Sprintf("%s%d", r.keyPref, id)
}

func (r *redisCache) Get(ctx context.Context, id int64) (*entity.Song, error) {
	key := fmt.Sprintf("song:%d", id)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var s entity.Song
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *redisCache) Set(ctx context.Context, s *entity.Song) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("song:%d", s.ID)
	log.Println(">>> SET REDIS:", key) // <- добавь это

	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		log.Println(">>> REDIS SET ERROR:", err)
		return err
	}

	err = r.client.Set(ctx, key, data, r.ttl).Err()
	return err
}

func (r *redisCache) Delete(ctx context.Context, id int64) error {
	return r.client.Del(ctx, r.key(id)).Err()
}
