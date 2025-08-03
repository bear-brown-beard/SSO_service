package containers

import (
	"database/sql"
	"sso/internal/logger"
	"sso/internal/repository"

	"github.com/go-redis/redis/v8"
)

type RepositoryContainer struct {
	TestAccountRepo repository.TestAccountRepository
	UserRepo        repository.UserRepository
	TokenRepo       repository.TokenRepository
	UserMindBoxRepo repository.UserMindBoxRepository
	logger          *logger.Logger
}

func NewRepositoryContainer(db *sql.DB, redisClient *redis.Client, logger *logger.Logger) (*RepositoryContainer, error) {
	container := &RepositoryContainer{
		logger: logger,
	}

	container.TestAccountRepo = repository.NewTestAccountRepository(db)
	container.UserRepo = repository.NewUserRepository(db)
	container.TokenRepo = repository.NewTokenRepository(db)
	container.UserMindBoxRepo = repository.NewUserMindBoxRepository(db)

	logger.Debug("All repositories initialized successfully")
	return container, nil
}

func (c *RepositoryContainer) GetTestAccountRepository() repository.TestAccountRepository {
	return c.TestAccountRepo
}

func (c *RepositoryContainer) GetUserRepository() repository.UserRepository {
	return c.UserRepo
}

func (c *RepositoryContainer) GetTokenRepository() repository.TokenRepository {
	return c.TokenRepo
}

func (c *RepositoryContainer) GetUserMindBoxRepository() repository.UserMindBoxRepository {
	return c.UserMindBoxRepo
}
