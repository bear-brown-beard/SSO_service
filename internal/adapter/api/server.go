package api

import (
	"net/http"
	handler "sso/internal/adapter/api/handler"
	"sso/internal/config"
	"sso/internal/logger"
)

func StartServer(handler *handler.VerificationHandler, cfg *config.Config, logger *logger.Logger) {
	router := NewRouter(handler, cfg)

	port := cfg.ServerPort
	if port == "" {
		port = "4053"
	}

	logger.Info("Server started", "port", port)
	logger.Error("Server stopped", "error", http.ListenAndServe(":"+port, router))
}
