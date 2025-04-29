package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	DB  DatabaseConfig
	App AppConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Schema   string
	Params   string
}

type AppConfig struct {
	LogLevel string
}

func (config *DatabaseConfig) ToConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		config.Username, config.Password, config.Host, config.Port, config.Database, config.Params)
}

func (config *AppConfig) ToSlogLevel() (slog.Level, error) {
	level := strings.ToLower(config.LogLevel)
	if level == "debug" {
		return slog.LevelDebug, nil
	} else if level == "info" {
		return slog.LevelInfo, nil
	} else if level == "warn" {
		return slog.LevelWarn, nil
	} else if level == "error" {
		return slog.LevelError, nil
	} else {
		return slog.LevelInfo, fmt.Errorf("invalid level: %s. Please use one of: debug, info, warn, error", level)
	}
}

func GetEnv(key string, required bool, missedEnvs *[]string) string {
	value, ok := os.LookupEnv(key)
	if !ok && required {
		*missedEnvs = append(*missedEnvs, key)
	}
	return value
}

func LoadConfig() (*Config, error) {
	missedEnvs := make([]string, 0)
	config := &Config{
		DB: DatabaseConfig{
			Host:     GetEnv("DB_HOST", true, &missedEnvs),
			Port:     GetEnv("DB_PORT", true, &missedEnvs),
			Username: GetEnv("DB_USERNAME", true, &missedEnvs),
			Password: GetEnv("DB_PASSWORD", true, &missedEnvs),
			Database: GetEnv("DB_DATABASE", true, &missedEnvs),
			Schema:   GetEnv("DB_SCHEMA", true, &missedEnvs),
		},
		App: AppConfig{
			LogLevel: GetEnv("LOG_LEVEL", false, &missedEnvs),
		},
	}
	var err error
	if len(missedEnvs) != 0 {
		err = fmt.Errorf("missing required environment variables: %s", strings.Join(missedEnvs, ","))
	}
	return config, err
}
