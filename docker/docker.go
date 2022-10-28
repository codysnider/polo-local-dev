package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/poloniex/polo-local-dev/output"
	"time"
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

func ContainerHealthCheckStream(ctx context.Context, container *types.Container, closer <-chan bool, outputChannel chan<- *types.HealthcheckResult) {

	if !ContainerHasHealthCheck(ctx, container) {
		return
	}

	containerInspect, inspectErr := dockerClient.ContainerInspect(ctx, container.ID)
	if inspectErr != nil {
		output.Error(inspectErr.Error())
		return
	}

	logsSent := map[time.Time]bool{}
	for _, logEntry := range containerInspect.State.Health.Log {
		if _, alreadySent := logsSent[logEntry.Start]; !alreadySent {
			outputChannel <- logEntry
			logsSent[logEntry.Start] = true
		}
	}

	outputWriteTick := time.NewTicker(time.Millisecond)
	for {
		select {
		case <-closer:
			outputWriteTick.Stop()
			return
		case <-outputWriteTick.C:
			for _, logEntry := range containerInspect.State.Health.Log {
				if _, alreadySent := logsSent[logEntry.Start]; !alreadySent {
					outputChannel <- logEntry
					logsSent[logEntry.Start] = true
				}
			}
		}
	}

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
