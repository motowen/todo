package container

import (
	"context"
	"fmt"
	"os"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	testcontainers.Container
	URI string
}

func SetupRedis(ctx context.Context) (*RedisContainer, error) {

	req := testcontainers.ContainerRequest{
		Image:        "redis:6",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}

	container, genericContainerErr := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if genericContainerErr != nil {
		return nil, genericContainerErr
	}

	mappedPort, mappedPortErr := container.MappedPort(ctx, "6379")
	if mappedPortErr != nil {
		return nil, mappedPortErr
	}

	hostIP, hostErr := container.Host(ctx)
	if hostErr != nil {
		return nil, hostErr
	}

	var endpoint string = fmt.Sprintf("%s:%s", hostIP, mappedPort.Port())
	os.Setenv("REDIS_ENDPOINT_LIST", endpoint)

	uri := fmt.Sprintf("redis://%s:%s", hostIP, mappedPort.Port())
	return &RedisContainer{Container: container, URI: uri}, nil
}
