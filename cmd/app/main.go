package main

import (
	"crud/cmd/app/config"
	"crud/cmd/app/config/log"
	"crud/cmd/app/server"
	"log/slog"
)

func main() {
	logger, logLevel := log.CreateLogger()

	appConfig, err := config.LoadConfig()
	if err != nil {
		logger.Error("Error loading config", slog.String("error", err.Error()))
		return
	}

	engine, dbPool, err := server.ConfigureAppEngine(appConfig, logLevel)
	if err != nil {
		logger.Error("Unable to configure app engine", slog.String("error", err.Error()))
		return
	}
	defer dbPool.Close()

	err = engine.Run(":8080")
	if err != nil {
		logger.Error("Error starting server", slog.String("error", err.Error()))
	}
}
