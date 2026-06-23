package service

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/containerd/errdefs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/moby/moby/api/pkg/stdcopy"
	moby_client "github.com/moby/moby/client"
	"github.com/rs/zerolog"
	container_repository "github.com/server-selfish/backend/internal/domain/repository/container"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

type (
	ContainerService interface {
		GetContainerStatus(ctx context.Context, name string, ui pgtype.UUID) (schema.ContainerStatusResponse, error)
		PauseContainer(ctx context.Context, name string, ui pgtype.UUID) error
		UnPauseContainer(ctx context.Context, name string, ui pgtype.UUID) error
		StopContainer(ctx context.Context, name string, ui pgtype.UUID) error
		StartContainer(ctx context.Context, name string, ui pgtype.UUID) error
		RestartContainer(ctx context.Context, name string, ui pgtype.UUID) error
		StreamLogs(ctx context.Context, userID pgtype.UUID, containerName string) (<-chan schema.ContainerLogEvent, <-chan error, error)
	}
	containerService struct {
		dc     *moby_client.Client
		cr     *container_repository.Queries
		conRep container_repository.ContainerRepository
		logger zerolog.Logger
	}
)

func NewContainerService(dc *moby_client.Client, cr *container_repository.Queries, conRep container_repository.ContainerRepository, logger zerolog.Logger) ContainerService {
	return &containerService{
		dc:     dc,
		cr:     cr,
		conRep: conRep,
		logger: logger,
	}
}

// StreamLogs implements [ContainerService].
func (c *containerService) StreamLogs(ctx context.Context, userID pgtype.UUID, containerName string) (<-chan schema.ContainerLogEvent, <-chan error, error) {

	events := make(chan schema.ContainerLogEvent, 1000)
	errs := make(chan error, 1)

	if _, err := c.cr.GetContainerByName(ctx, container_repository.GetContainerByNameParams{
		UserID: userID,
		Name:   containerName,
	}); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, defined_error.ErrContainerNotFound
		}
		return nil, nil, err
	}
	reader, err := c.conRep.GetContainerLogs(ctx, containerName)

	if err != nil {
		select {
		case errs <- err:
		default:
		}
		close(events)
		return events, errs, nil
	}

	go func() {
		defer func() {
			if err := reader.Close(); err != nil {
				c.logger.Error().Err(err).Msg("failed to close container log reader")
			}
		}()

		stdoutR, stdoutW := io.Pipe()
		stderrR, stderrW := io.Pipe()

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			defer wg.Done()

			defer func() {
				if err := stdoutW.Close(); err != nil {
					c.logger.Error().Err(err).Msg("failed to close stdout pipe writer")
				}
			}()
			defer func() {
				if err := stderrW.Close(); err != nil {
					c.logger.Error().Err(err).Msg("failed to close stderr pipe writer")
				}
			}()

			_, err := stdcopy.StdCopy(
				stdoutW,
				stderrW,
				reader,
			)

			if err != nil {
				select {
				case errs <- err:
				default:
				}
			}
		}()

		go func() {
			defer wg.Done()

			pkg.ScanStream(
				ctx,
				stdoutR,
				"stdout",
				events,
			)
		}()

		go func() {
			defer wg.Done()

			pkg.ScanStream(
				ctx,
				stderrR,
				"stderr",
				events,
			)
		}()

		<-ctx.Done()

		_ = stdoutR.Close()
		_ = stderrR.Close()

		wg.Wait()
		close(events)
		close(errs)
	}()

	return events, errs, nil
}

// RestartContainer implements [ContainerService].
func (c *containerService) RestartContainer(ctx context.Context, name string, ui pgtype.UUID) error {
	if _, err := c.cr.GetContainerByName(ctx, container_repository.GetContainerByNameParams{
		UserID: ui,
		Name:   name,
	}); err != nil {
		if err == pgx.ErrNoRows {
			return defined_error.ErrContainerNotFound
		}
		return err
	}
	if _, err := c.dc.ContainerRestart(ctx, name, moby_client.ContainerRestartOptions{}); err != nil {
		return err
	}
	return nil
}

// PauseContainer implements [ContainerService].
func (c *containerService) PauseContainer(ctx context.Context, name string, ui pgtype.UUID) error {
	if _, err := c.cr.GetContainerByName(ctx, container_repository.GetContainerByNameParams{
		UserID: ui,
		Name:   name,
	}); err != nil {
		if err == pgx.ErrNoRows {
			return defined_error.ErrContainerNotFound
		}
		return err
	}
	if _, err := c.dc.ContainerPause(ctx, name, moby_client.ContainerPauseOptions{}); err != nil {
		return err
	}
	return nil
}

// UnPauseContainer implements [ContainerService].
func (c *containerService) UnPauseContainer(ctx context.Context, name string, ui pgtype.UUID) error {
	if _, err := c.cr.GetContainerByName(ctx, container_repository.GetContainerByNameParams{
		UserID: ui,
		Name:   name,
	}); err != nil {
		if err == pgx.ErrNoRows {
			return defined_error.ErrContainerNotFound
		}
		return err
	}
	if _, err := c.dc.ContainerUnpause(ctx, name, moby_client.ContainerUnpauseOptions{}); err != nil {
		return err
	}
	return nil
}

// StartContainer implements [ContainerService].
func (c *containerService) StartContainer(ctx context.Context, name string, ui pgtype.UUID) error {
	if _, err := c.cr.GetContainerByName(ctx, container_repository.GetContainerByNameParams{
		UserID: ui,
		Name:   name,
	}); err != nil {
		if err == pgx.ErrNoRows {
			return defined_error.ErrContainerNotFound
		}
		return err
	}
	if _, err := c.dc.ContainerStart(ctx, name, moby_client.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

// StopContainer implements [ContainerService].
func (c *containerService) StopContainer(ctx context.Context, name string, ui pgtype.UUID) error {
	if _, err := c.cr.GetContainerByName(ctx, container_repository.GetContainerByNameParams{
		UserID: ui,
		Name:   name,
	}); err != nil {
		if err == pgx.ErrNoRows {
			return defined_error.ErrContainerNotFound
		}
		return err
	}
	if _, err := c.dc.ContainerStop(ctx, name, moby_client.ContainerStopOptions{}); err != nil {
		return err
	}
	return nil
}

// GetContainerStatus implements [ContainerService].
func (c *containerService) GetContainerStatus(ctx context.Context, name string, ui pgtype.UUID) (schema.ContainerStatusResponse, error) {
	if _, err := c.cr.GetContainerByName(ctx, container_repository.GetContainerByNameParams{
		UserID: ui,
		Name:   name,
	}); err != nil {
		if err == pgx.ErrNoRows {
			return schema.ContainerStatusResponse{}, defined_error.ErrContainerNotFound
		}
		return schema.ContainerStatusResponse{}, err
	}
	ir, err := c.dc.ContainerInspect(ctx, name, moby_client.ContainerInspectOptions{})
	if err != nil {
		if errdefs.IsNotFound(err) {
			return schema.ContainerStatusResponse{}, fmt.Errorf("%s: %s", defined_error.ErrContainerNotFound.Error(), name)
		}
		return schema.ContainerStatusResponse{}, err
	}
	return schema.ContainerStatusResponse{
		ContainerStatus: string(ir.Container.State.Status),
	}, nil
}
