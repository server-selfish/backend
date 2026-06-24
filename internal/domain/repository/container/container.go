package container_repository

import (
	"context"
	"io"

	moby_client "github.com/moby/moby/client"
	docker_infra "github.com/server-selfish/backend/internal/infra/docker"
)

type (
	ContainerRepository interface {
		GetContainerLogs(ctx context.Context, containerID string) (io.ReadCloser, error)
		ContainerRestart(ctx context.Context, container string, options moby_client.ContainerRestartOptions) (moby_client.ContainerRestartResult, error)
		ContainerPause(ctx context.Context, container string, options moby_client.ContainerPauseOptions) (moby_client.ContainerPauseResult, error)
		ContainerUnpause(ctx context.Context, container string, options moby_client.ContainerUnpauseOptions) (moby_client.ContainerUnpauseResult, error)
		ContainerStart(ctx context.Context, container string, options moby_client.ContainerStartOptions) (moby_client.ContainerStartResult, error)
		ContainerStop(ctx context.Context, container string, options moby_client.ContainerStopOptions) (moby_client.ContainerStopResult, error)
		ContainerInspect(ctx context.Context, container string, options moby_client.ContainerInspectOptions) (moby_client.ContainerInspectResult, error)
		ContainerRemove(ctx context.Context, container string, options moby_client.ContainerRemoveOptions) (moby_client.ContainerRemoveResult, error)
		ContainerCreate(ctx context.Context, options moby_client.ContainerCreateOptions) (moby_client.ContainerCreateResult, error)
		NetworkList(ctx context.Context, options moby_client.NetworkListOptions) (moby_client.NetworkListResult, error)
		NetworkCreate(ctx context.Context, name string, options moby_client.NetworkCreateOptions) (moby_client.NetworkCreateResult, error)
		ImageBuild(ctx context.Context, context io.Reader, options moby_client.ImageBuildOptions) (moby_client.ImageBuildResult, error)
		ImagePrune(ctx context.Context, opts moby_client.ImagePruneOptions) (moby_client.ImagePruneResult, error)
		ImageRemove(ctx context.Context, image string, options moby_client.ImageRemoveOptions) (moby_client.ImageRemoveResult, error)
	}
	containerRepository struct {
		dc docker_infra.DockerInfra
	}
)

func NewContainerRepository(dc docker_infra.DockerInfra) ContainerRepository {
	return containerRepository{
		dc: dc,
	}
}

// ContainerCreate implements [ContainerRepository].
func (c containerRepository) ContainerCreate(ctx context.Context, options moby_client.ContainerCreateOptions) (moby_client.ContainerCreateResult, error) {
	return c.dc.ContainerCreate(ctx, options)
}

// ContainerRemove implements [ContainerRepository].
func (c containerRepository) ContainerRemove(ctx context.Context, container string, options moby_client.ContainerRemoveOptions) (moby_client.ContainerRemoveResult, error) {
	return c.dc.ContainerRemove(ctx, container, options)
}

// ImageBuild implements [ContainerRepository].
func (c containerRepository) ImageBuild(ctx context.Context, context io.Reader, options moby_client.ImageBuildOptions) (moby_client.ImageBuildResult, error) {
	return c.dc.ImageBuild(ctx, context, options)
}

// ImagePrune implements [ContainerRepository].
func (c containerRepository) ImagePrune(ctx context.Context, opts moby_client.ImagePruneOptions) (moby_client.ImagePruneResult, error) {
	return c.dc.ImagePrune(ctx, opts)
}

// ImageRemove implements [ContainerRepository].
func (c containerRepository) ImageRemove(ctx context.Context, image string, options moby_client.ImageRemoveOptions) (moby_client.ImageRemoveResult, error) {
	return c.dc.ImageRemove(ctx, image, options)
}

// NetworkCreate implements [ContainerRepository].
func (c containerRepository) NetworkCreate(ctx context.Context, name string, options moby_client.NetworkCreateOptions) (moby_client.NetworkCreateResult, error) {
	return c.dc.NetworkCreate(ctx, name, options)
}

// NetworkList implements [ContainerRepository].
func (c containerRepository) NetworkList(ctx context.Context, options moby_client.NetworkListOptions) (moby_client.NetworkListResult, error) {
	return c.dc.NetworkList(ctx, options)
}

// ContainerInspect implements [ContainerRepository].
func (c containerRepository) ContainerInspect(ctx context.Context, container string, options moby_client.ContainerInspectOptions) (moby_client.ContainerInspectResult, error) {
	return c.dc.ContainerInspect(ctx, container, options)
}

// ContainerPause implements [ContainerRepository].
func (c containerRepository) ContainerPause(ctx context.Context, container string, options moby_client.ContainerPauseOptions) (moby_client.ContainerPauseResult, error) {
	return c.dc.ContainerPause(ctx, container, options)
}

// ContainerRestart implements [ContainerRepository].
func (c containerRepository) ContainerRestart(ctx context.Context, container string, options moby_client.ContainerRestartOptions) (moby_client.ContainerRestartResult, error) {
	return c.dc.ContainerRestart(ctx, container, options)
}

// ContainerStart implements [ContainerRepository].
func (c containerRepository) ContainerStart(ctx context.Context, container string, options moby_client.ContainerStartOptions) (moby_client.ContainerStartResult, error) {
	return c.dc.ContainerStart(ctx, container, options)
}

// ContainerStop implements [ContainerRepository].
func (c containerRepository) ContainerStop(ctx context.Context, container string, options moby_client.ContainerStopOptions) (moby_client.ContainerStopResult, error) {
	return c.dc.ContainerStop(ctx, container, options)
}

// ContainerUnpause implements [ContainerRepository].
func (c containerRepository) ContainerUnpause(ctx context.Context, container string, options moby_client.ContainerUnpauseOptions) (moby_client.ContainerUnpauseResult, error) {
	return c.dc.ContainerUnpause(ctx, container, options)
}

// GetContainerLogs implements [ContainerRepository].
func (c containerRepository) GetContainerLogs(ctx context.Context, containerName string) (io.ReadCloser, error) {
	return c.dc.ContainerLogs(ctx, containerName, moby_client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "600",
		// Details:    false,
	})
}
