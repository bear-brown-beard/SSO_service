package containers

import (
	"sso/internal/logger"
	"sso/internal/service"

	"github.com/go-redis/redis/v8"
)

type CacheContainer struct {
	codeCache service.CacheService
	logger    *logger.Logger
}

func NewCacheContainer(redisClient *redis.Client, logger *logger.Logger) (*CacheContainer, error) {
	container := &CacheContainer{
		logger: logger,
	}

	container.codeCache = service.NewRedisCache(redisClient)

	logger.Debug("All cache services initialized successfully")
	return container, nil
}

func (c *CacheContainer) GetCodeCache() service.CacheService {
	return c.codeCache
}
