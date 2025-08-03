package containers

import (
	"context"
	"database/sql"
	"fmt"
	"sso/internal/config"
	"sso/internal/logger"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	maxAttempts       = 5
	connectionTimeout = 5 * time.Second
	retryDelay        = 3 * time.Second
)

type DatabaseContainer struct {
	DB     *sql.DB
	logger *logger.Logger
}

func NewDatabaseContainer(cfg *config.Config, logger *logger.Logger) (container *DatabaseContainer, err error) {
	container = &DatabaseContainer{
		logger: logger,
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.Mysql.DBUser,
		cfg.Mysql.DBPassword,
		cfg.Mysql.DBHost,
		cfg.Mysql.DBPort,
		cfg.Mysql.DBName)

	var db *sql.DB
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		logger.Debug("Attempting to connect to database", "attempt", attempt, "max_attempts", maxAttempts)

		db, err = sql.Open("mysql", dsn)
		if err != nil {
			logger.Error("Failed to open database connection", "error", err, "attempt", attempt)
			time.Sleep(retryDelay)
			continue
		}

		db.SetMaxOpenConns(cfg.Mysql.DBMaxOpenConns)
		db.SetMaxIdleConns(cfg.Mysql.DBMaxIdleConns)
		db.SetConnMaxLifetime(cfg.Mysql.DBConnMaxLifetime)

		ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
		err = db.PingContext(ctx)
		cancel()

		if err == nil {
			logger.Info("Database connected successfully", "host", cfg.Mysql.DBHost, "database", cfg.Mysql.DBName)
			container.DB = db
			return container, nil
		}

		logger.Error("Failed to ping database", "error", err, "attempt", attempt)
		db.Close()
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts", maxAttempts)
}

func (c *DatabaseContainer) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
