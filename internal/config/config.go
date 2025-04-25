package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DB DatabaseConfig
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

func (config *DatabaseConfig) ToConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		config.Username, config.Password, config.Host, config.Port, config.Database, config.Params)
}

func GetEnv(key string, missedEnvs *[]string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		*missedEnvs = append(*missedEnvs, key)
	}
	return value
}

func LoadConfig() (*Config, error) {
	missedEnvs := make([]string, 0)
	config := &Config{DB: DatabaseConfig{
		Host:     GetEnv("DB_HOST", &missedEnvs),
		Port:     GetEnv("DB_PORT", &missedEnvs),
		Username: GetEnv("DB_USERNAME", &missedEnvs),
		Password: GetEnv("DB_PASSWORD", &missedEnvs),
		Database: GetEnv("DB_DATABASE", &missedEnvs),
		Schema:   GetEnv("DB_SCHEMA", &missedEnvs),
	}}
	var err error
	if len(missedEnvs) != 0 {
		err = fmt.Errorf("missing required environment variables: %s", strings.Join(missedEnvs, ","))
	}
	return config, err
}
