package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	moby_client "github.com/moby/moby/client"
	container_repository "github.com/server-selfish/backend/internal/domain/repository/container"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/pkg"
)

type (
	ContainerService interface {
		GetContainerStatus(ctx context.Context, name string, ui pgtype.UUID) (schema.ContainerStatusResponse, error)
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

// GetContainerStatus implements [ContainerService].
func (c *containerService) GetContainerStatus(ctx context.Context, name string, ui pgtype.UUID) (schema.ContainerStatusResponse, error) {
	if _, err := c.cr.GetContainerByName(ctx, container_repository.GetContainerByNameParams{
		UserID: ui,
		Name:   name,
	}); err != nil {
		if err == pgx.ErrNoRows {
			return schema.ContainerStatusResponse{}, pkg.ErrNotFound
		}
		return schema.ContainerStatusResponse{}, err
	}
	ir, err := c.dc.ContainerInspect(ctx, name, moby_client.ContainerInspectOptions{})
	if err != nil {
		return schema.ContainerStatusResponse{}, err
	}
	return schema.ContainerStatusResponse{
		ContainerStatus: string(ir.Container.State.Status),
	}, nil
}
