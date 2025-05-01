package server

import (
	"context"
	"crud/cmd/app/config"
	logConfig "crud/cmd/app/config/log"
	"crud/internal/model"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

const postgresTestPassword = "testpassword"

func TestIntegrationApp(t *testing.T) {
	_, logLevel := logConfig.CreateLogger()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:17-alpine",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env:          map[string]string{"POSTGRES_PASSWORD": postgresTestPassword},
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
		Password: postgresTestPassword,
		Database: "postgres",
		Schema:   "public",
		Params:   "",
	}, App: config.AppConfig{
		LogLevel: "info",
		AppMode:  "test",
	}}

	engine, dbPool, err := ConfigureAppEngine(&appConfig, logLevel)
	assert.NoError(t, err)
	server := httptest.NewServer(engine.Handler())
	client := server.Client()

	httpResponse, err := client.Get(server.URL + "/api/v1/user/")
	assert.NoError(t, err)
	body, err := io.ReadAll(httpResponse.Body)
	var usersResponse []model.UserResponse
	err = json.Unmarshal(body, &usersResponse)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(usersResponse))

	server.Close()
	dbPool.Close()
	testcontainers.CleanupContainer(t, postgres)
	require.NoError(t, err)
}
