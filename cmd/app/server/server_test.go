package server

import (
	"context"
	"crud/cmd/app/config"
	logConfig "crud/cmd/app/config/log"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"net/http"
	"strings"
	"testing"
)

const POSTGRES_TEST_PASSWORD = "testpassword"

//func startServer(t *testing.T) (testcontainers.Container, string) {
//
//}

func TestIntegrationApp(t *testing.T) {
	_, logLevel := logConfig.CreateLogger()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:17-alpine",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env:          map[string]string{"POSTGRES_PASSWORD": POSTGRES_TEST_PASSWORD},
	}
	ctx := context.Background()
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	endpoint, err := postgres.Endpoint(ctx, "")
	assert.NoError(t, err)
	hostAndPort := strings.Split(endpoint, ":")
	appConfig := config.Config{DB: config.DatabaseConfig{
		Host:     hostAndPort[0],
		Port:     hostAndPort[1],
		Username: "postgres",
		Password: POSTGRES_TEST_PASSWORD,
		Database: "postgres",
		Schema:   "public",
		Params:   "",
	}, App: config.AppConfig{
		LogLevel: "info",
		AppMode:  "test",
	}}

	engine, dbPool, err := ConfigureAppEngine(&appConfig, logLevel)
	assert.NoError(t, err)
	srv := http.Server{
		Addr:    ":8080",
		Handler: engine.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	err = srv.Shutdown(ctx)
	dbPool.Close()
	testcontainers.CleanupContainer(t, postgres)
	require.NoError(t, err)
}
