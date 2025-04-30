package main

import (
	"github.com/Cladkoewka/http-load-balancer/internal/config"
	"github.com/Cladkoewka/http-load-balancer/internal/logger"
	"github.com/Cladkoewka/http-load-balancer/internal/server"
)

func main() {
	logger.InitLogger()

	// Load configuration
	cfg, err := config.Load("../../configs/config.yaml")
	if err != nil {
		logger.Log.Error("failed to load config", "error", err)
	}

	if err := server.StartServer(cfg); err != nil {
		logger.Log.Error("Failed to start server:", "error", err)
	}
}
