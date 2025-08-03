package containers

import (
	"sso/internal/config"
	"sso/internal/logger"
)

type LoggerContainer struct {
	Logger *logger.Logger
}

func NewLoggerContainer(cfg *config.Config) *LoggerContainer {
	log := setupLogger(cfg)
	return &LoggerContainer{
		Logger: log,
	}
}

func setupLogger(cfg *config.Config) *logger.Logger {
	var log *logger.Logger
	switch cfg.AppEnv {
	case "local":
		log = logger.NewLogger(true)
	case "dev":
		log = logger.NewLogger(false)
	case "stage":
		log = logger.NewLogger(false)
	case "prod":
		log = logger.NewLogger(false)
	default:
		log = logger.NewLogger(true)
	}
	return log
}
