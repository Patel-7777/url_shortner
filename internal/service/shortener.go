package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/drashti/url_shortner/internal/storage"
)

type ShortenerService struct {
	postgres *storage.PostgresStorage
	redis    *storage.RedisStorage
}

func NewShortenerService(postgres *storage.PostgresStorage, redis *storage.RedisStorage) *ShortenerService {
	return &ShortenerService{
		postgres: postgres,
		redis:    redis,
	}
}

func (s *ShortenerService) GenerateShortCode() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:8], nil
}

func (s *ShortenerService) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	shortCode, err := s.GenerateShortCode()
	if err != nil {
		return "", fmt.Errorf("error generating short code: %v", err)
	}

	if err := s.postgres.CreateURL(shortCode, originalURL); err != nil {
		return "", fmt.Errorf("error storing URL: %v", err)
	}

	if err := s.redis.SetURL(ctx, shortCode, originalURL); err != nil {
		// Log the error but don't fail the request
		fmt.Printf("Error caching URL in Redis: %v\n", err)
	}

	return shortCode, nil
}

func (s *ShortenerService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	// Try Redis first
	url, err := s.redis.GetURL(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("error getting URL from Redis: %v", err)
	}
	if url != "" {
		return url, nil
	}

	// Fall back to PostgreSQL
	record, err := s.postgres.GetURL(shortCode)
	if err != nil {
		return "", fmt.Errorf("error getting URL from PostgreSQL: %v", err)
	}
	if record == nil {
		return "", nil
	}

	// Cache the result in Redis
	if err := s.redis.SetURL(ctx, shortCode, record.OriginalURL); err != nil {
		fmt.Printf("Error caching URL in Redis: %v\n", err)
	}

	// Increment visit count
	if err := s.postgres.IncrementVisitCount(shortCode); err != nil {
		fmt.Printf("Error incrementing visit count: %v\n", err)
	}

	return record.OriginalURL, nil
}

func (s *ShortenerService) Close() error {
	if err := s.postgres.Close(); err != nil {
		return err
	}
	return s.redis.Close()
}
