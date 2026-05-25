package docker_client

import (
	"context"

	"github.com/moby/moby/client"
	"github.com/rs/zerolog"
)

func NewDockerClient(logger zerolog.Logger) *client.Client {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initiate docker client")
	}
	if _, err := cli.Ping(context.Background(), client.PingOptions{}); err != nil {
		logger.Fatal().Err(err).Msg("failed to connect with docker")
	}
	return cli
}
