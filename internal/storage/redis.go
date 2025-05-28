package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	urlCacheTTL  = 24 * time.Hour
	urlKeyPrefix = "short:"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(host, port, password string, db int) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to redis: %v", err)
	}

	return &RedisStorage{client: client}, nil
}

func (s *RedisStorage) SetURL(ctx context.Context, shortCode, originalURL string) error {
	key := urlKeyPrefix + shortCode
	return s.client.Set(ctx, key, originalURL, urlCacheTTL).Err()
}

func (s *RedisStorage) GetURL(ctx context.Context, shortCode string) (string, error) {
	key := urlKeyPrefix + shortCode
	url, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *RedisStorage) Close() error {
	return s.client.Close()
}
