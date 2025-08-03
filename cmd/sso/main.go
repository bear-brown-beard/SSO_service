package main

import (
	"os"
	"sso/internal/adapter/api"
	apiHandler "sso/internal/adapter/api/handler"
	"sso/internal/config"
	container "sso/internal/di"
)

func main() {
	cfg := config.Load()

	container, err := container.NewContainer(cfg)
	if err != nil {
		os.Exit(1)
	}
	defer container.Close()

	SSOService := container.GetSSOService()
	logger := container.GetLogger()

	handler := apiHandler.NewVerificationHandler(SSOService, logger)
	api.StartServer(handler, cfg, logger)
}
