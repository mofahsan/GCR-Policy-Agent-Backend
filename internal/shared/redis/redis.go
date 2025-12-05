package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func Init(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	client = redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}

func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
