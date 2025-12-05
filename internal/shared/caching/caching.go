package caching

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheService defines the interface for a cache.
type CacheService interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

// RedisCacheService is the Redis implementation of the CacheService.
type RedisCacheService struct {
	client *redis.Client
}

// NewRedisCacheService creates a new RedisCacheService.
func NewRedisCacheService(client *redis.Client) *RedisCacheService {
	return &RedisCacheService{client: client}
}

// Get retrieves an item from the cache.
func (s *RedisCacheService) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Set adds an item to the cache.
func (s *RedisCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, b, expiration).Err()
}

// Delete removes an item from the cache.
func (s *RedisCacheService) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}
