package containers

import (
	"fmt"
	"sso/internal/config"
	"sso/internal/di/containers"
	"sso/internal/logger"
	"sso/internal/service"
)

type Container struct {
	config           *config.Config
	loggerContainer  *containers.LoggerContainer
	dbContainer      *containers.DatabaseContainer
	redisContainer   *containers.RedisContainer
	cacheContainer   *containers.CacheContainer
	repoContainer    *containers.RepositoryContainer
	serviceContainer *containers.ServiceContainer
}

func NewContainer(cfg *config.Config) (*Container, error) {
	container := &Container{
		config: cfg,
	}

	loggerContainer := containers.NewLoggerContainer(cfg)
	container.loggerContainer = loggerContainer

	dbContainer, err := containers.NewDatabaseContainer(cfg, loggerContainer.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database container: %w", err)
	}
	container.dbContainer = dbContainer

	redisContainer, err := containers.NewRedisContainer(cfg, loggerContainer.Logger)
	if err != nil {
		container.Close()
		return nil, fmt.Errorf("failed to create redis container: %w", err)
	}
	container.redisContainer = redisContainer

	cacheContainer, err := containers.NewCacheContainer(redisContainer.RedisClient, loggerContainer.Logger)
	if err != nil {
		container.Close()
		return nil, fmt.Errorf("failed to create cache container: %w", err)
	}
	container.cacheContainer = cacheContainer

	repoContainer, err := containers.NewRepositoryContainer(dbContainer.DB, redisContainer.RedisClient, loggerContainer.Logger)
	if err != nil {
		container.Close()
		return nil, fmt.Errorf("failed to create repository container: %w", err)
	}
	container.repoContainer = repoContainer

	serviceContainer, err := containers.NewServiceContainer(repoContainer, cacheContainer, cfg, loggerContainer.Logger)
	if err != nil {
		container.Close()
		return nil, fmt.Errorf("failed to create service container: %w", err)
	}
	container.serviceContainer = serviceContainer

	return container, nil
}

func (c *Container) GetSSOService() service.SSOService {
	return c.serviceContainer.GetSSOService()
}

func (c *Container) GetLogger() *logger.Logger {
	return c.loggerContainer.Logger
}

func (c *Container) Close() {
	if c.redisContainer != nil {
		if err := c.redisContainer.Close(); err != nil {
			c.loggerContainer.Logger.Error("Error closing Redis connection", "error", err)
		}
	}

	if c.dbContainer != nil {
		if err := c.dbContainer.Close(); err != nil {
			c.loggerContainer.Logger.Error("Error closing database connection", "error", err)
		}
	}
}
