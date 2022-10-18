package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/poloniex/polo-local-dev/output"
)

var (
	dockerClient *client.Client
)

func init() {

	var dockerClientErr error

	dockerClient, dockerClientErr = client.NewClientWithOpts(client.FromEnv)
	if dockerClientErr != nil {
		output.Error(dockerClientErr.Error())
		return
	}
}

func Containers(ctx context.Context) ([]types.Container, error) {
	return dockerClient.ContainerList(ctx, types.ContainerListOptions{})
}

func ContainerHasHealthCheck(ctx context.Context, container *types.Container) bool {
	containerInspect, inspectErr := dockerClient.ContainerInspect(ctx, container.ID)
	if inspectErr != nil {
		output.Error(inspectErr.Error())
		return false
	}

	return containerInspect.State.Health != nil
}

func ContainerIsHealthy(container *types.Container) bool {
	if ContainerHasHealthCheck(context.Background(), container) {
		containerInspect, inspectErr := dockerClient.ContainerInspect(context.Background(), container.ID)
		if inspectErr != nil {
			output.Error(inspectErr.Error())
			return false
		}

		return containerInspect.State.Health.Status == "healthy"
	}

	return false
}
