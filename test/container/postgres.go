package container

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
)

type PostgresContainer struct {
	testcontainers.Container
	URI string
}

func SetupPostgres(ctx context.Context) (*PostgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:      "postgres:14-alpine3.15",
		Entrypoint: nil,
		Env: map[string]string{
			"POSTGRES_DB":       "postgres",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "pgpassword",
		},
		ExposedPorts: []string{"5432/tcp"},
		Name:         "postgres",
		User:         "postgres",
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		AutoRemove:   true,
	}

	container, genericContainerErr := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if genericContainerErr != nil {
		return nil, genericContainerErr
	}

	mappedPort, mappedPortErr := container.MappedPort(ctx, "5432")
	if mappedPortErr != nil {
		return nil, mappedPortErr
	}

	hostIP, hostErr := container.Host(ctx)
	if hostErr != nil {
		return nil, hostErr
	}

	os.Setenv("DATABASE_NAME", "postgres")
	os.Setenv("DATABASE_USERNAME", "postgres")
	os.Setenv("DATABASE_PASSWORD", "pgpassword")
	os.Setenv("DATABASE_HOST", hostIP)
	os.Setenv("DATABASE_PORT", mappedPort.Port())

	uri := fmt.Sprintf("postgres://%s:%s@%s:%s", "postgres", "pgpassword", hostIP, mappedPort.Port())
	return &PostgresContainer{Container: container, URI: uri}, nil
}
