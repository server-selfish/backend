package service

import (
	"context"
	"fmt"

	"github.com/containerd/errdefs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	moby_client "github.com/moby/moby/client"
	container_repository "github.com/server-selfish/backend/internal/domain/repository/container"
	"github.com/server-selfish/backend/internal/domain/schema"
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
	}
	containerService struct {
		dc *moby_client.Client
		cr *container_repository.Queries
	}
)

func NewContainerService(dc *moby_client.Client, cr *container_repository.Queries) ContainerService {
	return &containerService{
		dc: dc,
		cr: cr,
	}
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
