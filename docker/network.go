package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/poloniex/polo-local-dev/output"
)

func CreateNetwork(name string) error {

	// Get all docker networks
	networks, networkListErr := dockerClient.NetworkList(context.Background(), types.NetworkListOptions{})
	if networkListErr != nil {
		return networkListErr
	}

	for _, network := range networks {
		output.Plain(network.Name)
	}

	return nil
}
