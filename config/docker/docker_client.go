package docker_client

import (
	"context"

	"github.com/moby/moby/client"
)

func NewDockerClient() *client.Client {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		panic("failed to initiate docker client")
	}
	if _, err := cli.Ping(context.Background(), client.PingOptions{}); err != nil {
		panic("failed to connect with docker")
	}
	return cli
}
