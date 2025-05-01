package main

import (
	"crud/cmd/app/config"
	"crud/cmd/app/config/log"
	"crud/cmd/app/server"
	"log/slog"
	"os"
)

func main() {
	logger, logLevel := log.CreateLogger()

	appConfig, err := config.LoadConfig()
	if err != nil {
		logger.Error("Error loading config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = server.Run(appConfig, logLevel)
	if err != nil {
		logger.Error("Something went wrong")
		return
	}
}
