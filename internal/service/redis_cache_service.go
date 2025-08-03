package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheService interface {
	SaveCode(ctx context.Context, phone, code string, ttl time.Duration) error
	GetCode(ctx context.Context, phone string) (string, error)
	DeleteCode(ctx context.Context, phone string) error

	AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	if client == nil {
		log.Fatalf("Provided Redis client is nil to NewRedisCache")
	}
	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) SaveCode(ctx context.Context, phone, code string, ttl time.Duration) error {
	err := r.client.Set(ctx, phone, code, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save code %s to Redis for phone %s: %w", code, phone, err)
	}
	return nil
}

func (r *RedisCache) GetCode(ctx context.Context, phone string) (string, error) {
	val, err := r.client.Get(ctx, phone).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get code for phone %s from Redis: %w", phone, err)
	}
	return val, nil
}

func (r *RedisCache) DeleteCode(ctx context.Context, phone string) error {
	err := r.client.Del(ctx, phone).Err()
	if err != nil {
		return fmt.Errorf("failed to delete code for phone %s from Redis: %w", phone, err)
	}
	return nil
}

func (r *RedisCache) AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	key := fmt.Sprintf("blacklist_token:%s", token)
	err := r.client.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}
	return nil
}

func (r *RedisCache) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist_token:%s", token)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check token in blacklist: %w", err)
	}
	return val == "1", nil
}

var _ CacheService = (*RedisCache)(nil)
