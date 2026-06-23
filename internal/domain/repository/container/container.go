package container_repository

import (
	"context"
	"io"

	moby_client "github.com/moby/moby/client"
)

type (
	ContainerRepository interface {
		GetContainerLogs(
			ctx context.Context,
			containerID string,
		) (io.ReadCloser, error)
	}
	containerRepository struct {
		dc *moby_client.Client
	}
)

func NewContainerRepository(dc *moby_client.Client) ContainerRepository {
	return containerRepository{
		dc: dc,
	}
}

// GetContainerLogs implements [ContainerRepository].
func (c containerRepository) GetContainerLogs(ctx context.Context, containerName string) (io.ReadCloser, error) {
	return c.dc.ContainerLogs(ctx, containerName, moby_client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		// Since:      "",
		// Until:      "",
		// Timestamps: false,
		Follow: true,
		Tail:   "600",
		// Details:    false,
	})
}
