package containers

import (
	"fmt"
	"sso/internal/config"
	"sso/internal/logger"

	"github.com/go-redis/redis/v8"
)

type RedisContainer struct {
	RedisClient *redis.Client
	logger      *logger.Logger
}

func NewRedisContainer(cfg *config.Config, logger *logger.Logger) (*RedisContainer, error) {
	container := &RedisContainer{
		logger: logger,
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	container.RedisClient = client
	return container, nil
}

func (c *RedisContainer) Close() error {
	if c.RedisClient != nil {
		return c.RedisClient.Close()
	}
	return nil
}
