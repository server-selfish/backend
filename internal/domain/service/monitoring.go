package service

import (
	"context"
	"time"

	"github.com/server-selfish/backend/internal/constant"
	monitoring_repository "github.com/server-selfish/backend/internal/domain/repository/monitoring"
	"github.com/server-selfish/backend/internal/domain/schema"
)

type (
	PrometheusService interface {
		GetCPUUsage(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetMemoryUsage(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetNetworkTx(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetNetworkRx(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetIORead(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetIOWrite(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
	}
	prometheusService struct {
		pr monitoring_repository.PrometheusRepository
	}
)

func NewPrometheusService(pr monitoring_repository.PrometheusRepository) PrometheusService {
	return prometheusService{
		pr: pr,
	}
}

// GetCPUUsage implements [PrometheusService].
func (p prometheusService) GetCPUUsage(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	duration := params.EndTime.Sub(params.StartTime)

	params.Step = duration / constant.MAX_POINTS

	if params.Step < 5*time.Second {
		params.Step = 5 * time.Second
	}
	return p.pr.GetCPUUsage(ctx, params)
}

// GetIORead implements [PrometheusService].
func (p prometheusService) GetIORead(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	duration := params.EndTime.Sub(params.StartTime)

	params.Step = duration / constant.MAX_POINTS

	if params.Step < 5*time.Second {
		params.Step = 5 * time.Second
	}
	return p.pr.GetIORead(ctx, params)
}

// GetIOWrite implements [PrometheusService].
func (p prometheusService) GetIOWrite(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	duration := params.EndTime.Sub(params.StartTime)

	params.Step = duration / constant.MAX_POINTS

	if params.Step < 5*time.Second {
		params.Step = 5 * time.Second
	}
	return p.pr.GetIOWrite(ctx, params)
}

// GetMemoryUsage implements [PrometheusService].
func (p prometheusService) GetMemoryUsage(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	duration := params.EndTime.Sub(params.StartTime)

	params.Step = duration / constant.MAX_POINTS

	if params.Step < 5*time.Second {
		params.Step = 5 * time.Second
	}
	return p.pr.GetMemoryUsage(ctx, params)
}

// GetNetworkRx implements [PrometheusService].
func (p prometheusService) GetNetworkRx(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	duration := params.EndTime.Sub(params.StartTime)

	params.Step = duration / constant.MAX_POINTS

	if params.Step < 5*time.Second {
		params.Step = 5 * time.Second
	}
	return p.pr.GetNetworkRx(ctx, params)
}

// GetNetworkTx implements [PrometheusService].
func (p prometheusService) GetNetworkTx(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	duration := params.EndTime.Sub(params.StartTime)

	params.Step = duration / constant.MAX_POINTS

	if params.Step < 5*time.Second {
		params.Step = 5 * time.Second
	}
	return p.pr.GetNetworkTx(ctx, params)
}
