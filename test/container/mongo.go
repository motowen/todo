package container

import (
	"context"
	"fmt"
	"os"

	"github.com/testcontainers/testcontainers-go"
)

type mongoContainer struct {
	testcontainers.Container
	URI string
}

func SetupMongo(ctx context.Context) (*mongoContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:3.6",
		Entrypoint:   nil,
		ExposedPorts: []string{"27017/tcp"},
		Name:         "mongo_test",
		AutoRemove:   true,
	}

	container, genericContainerErr := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if genericContainerErr != nil {
		return nil, genericContainerErr
	}

	mappedPort, mappedPortErr := container.MappedPort(ctx, "27017/tcp")
	if mappedPortErr != nil {
		return nil, mappedPortErr
	}

	hostIP, hostErr := container.Host(ctx)
	if hostErr != nil {
		return nil, hostErr
	}

	os.Setenv("DATABASE_HOST", hostIP)
	os.Setenv("DATABASE_PORT", mappedPort.Port())

	uri := fmt.Sprintf("mongodb://%s:%s", hostIP, mappedPort.Port())

	return &mongoContainer{Container: container, URI: uri}, nil
}
