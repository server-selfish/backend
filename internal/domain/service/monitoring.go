package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/server-selfish/backend/internal/constant"
	container_repository "github.com/server-selfish/backend/internal/domain/repository/container"
	monitoring_repository "github.com/server-selfish/backend/internal/domain/repository/monitoring"
	"github.com/server-selfish/backend/internal/domain/schema"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

type (
	PrometheusService interface {
		GetCPUUsage(ctx context.Context, userID pgtype.UUID, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetMemoryUsage(ctx context.Context, userID pgtype.UUID, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetNetworkTx(ctx context.Context, userID pgtype.UUID, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetNetworkRx(ctx context.Context, userID pgtype.UUID, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetIORead(ctx context.Context, userID pgtype.UUID, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetIOWrite(ctx context.Context, userID pgtype.UUID, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
	}
	prometheusService struct {
		pr monitoring_repository.PrometheusRepository
		cr *container_repository.Queries
	}
)

func NewPrometheusService(pr monitoring_repository.PrometheusRepository, cr *container_repository.Queries) PrometheusService {
	return prometheusService{
		pr: pr,
		cr: cr,
	}
}

// GetCPUUsage implements [PrometheusService].
func (p prometheusService) GetCPUUsage(ctx context.Context, userID pgtype.UUID, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	if err := p.ensureOwnership(ctx, userID, params.ContainerName); err != nil {
		return nil, err
	}
	duration := params.EndTime.Sub(params.StartTime)

	params.Step = duration / constant.MAX_POINTS

	if params.Step < 5*time.Second {
		params.Step = 5 * time.Second
	}
	return p.pr.GetCPUUsage(ctx, params)
}

// GetIORead implements [PrometheusService].
func (p prometheusService) GetIORead(ctx context.Context, userID pgtype.UUID, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	if err := p.ensureOwnership(ctx, userID, params.ContainerName); err != nil {
		return nil, err
	}
	duration := params.EndTime.Sub(params.StartTime)

	params.Step = duration / constant.MAX_POINTS

	if params.Step < 5*time.Second {
		params.Step = 5 * time.Second
	}
	return p.pr.GetIORead(ctx, params)
}

// GetIOWrite implements [PrometheusService].
func (p prometheusService) GetIOWrite(ctx context.Context, userID pgtype.UUID, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	if err := p.ensureOwnership(ctx, userID, params.ContainerName); err != nil {
		return nil, err
	}
	duration := params.EndTime.Sub(params.StartTime)

	params.Step = duration / constant.MAX_POINTS

	if params.Step < 5*time.Second {
		params.Step = 5 * time.Second
	}
	return p.pr.GetIOWrite(ctx, params)
}

// GetMemoryUsage implements [PrometheusService].
func (p prometheusService) GetMemoryUsage(ctx context.Context, userID pgtype.UUID, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	if err := p.ensureOwnership(ctx, userID, params.ContainerName); err != nil {
		return nil, err
	}
	duration := params.EndTime.Sub(params.StartTime)

	params.Step = duration / constant.MAX_POINTS

	if params.Step < 5*time.Second {
		params.Step = 5 * time.Second
	}
	return p.pr.GetMemoryUsage(ctx, params)
}

// GetNetworkRx implements [PrometheusService].
func (p prometheusService) GetNetworkRx(ctx context.Context, userID pgtype.UUID, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	if err := p.ensureOwnership(ctx, userID, params.ContainerName); err != nil {
		return nil, err
	}
	duration := params.EndTime.Sub(params.StartTime)

	params.Step = duration / constant.MAX_POINTS

	if params.Step < 5*time.Second {
		params.Step = 5 * time.Second
	}
	return p.pr.GetNetworkRx(ctx, params)
}

// GetNetworkTx implements [PrometheusService].
func (p prometheusService) GetNetworkTx(ctx context.Context, userID pgtype.UUID, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	if err := p.ensureOwnership(ctx, userID, params.ContainerName); err != nil {
		return nil, err
	}
	duration := params.EndTime.Sub(params.StartTime)

	params.Step = duration / constant.MAX_POINTS

	if params.Step < 5*time.Second {
		params.Step = 5 * time.Second
	}
	return p.pr.GetNetworkTx(ctx, params)
}

func (p prometheusService) ensureOwnership(ctx context.Context, userID pgtype.UUID, containerName string) error {
	if _, err := p.cr.GetContainerByName(ctx, container_repository.GetContainerByNameParams{
		UserID: userID,
		Name:   containerName,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return defined_error.ErrContainerNotFound
		}
		return err
	}
	return nil
}
