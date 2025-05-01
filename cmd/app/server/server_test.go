package server

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func TestIntegrationApp(t *testing.T) {
	//_, logLevel := log.CreateLogger()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:17-alpine",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env:          map[string]string{"POSTGRES_PASSWORD": "testpassword"},
	}
	postgres, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		println(err.Error())
	}
	//appConfig := config.Config{DB: {}, App: {}}

	//Run(&appConfig, logLevel)

	testcontainers.CleanupContainer(t, postgres)
	require.NoError(t, err)
}
