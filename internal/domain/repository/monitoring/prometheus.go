package monitoring_repository

import (
	"context"
	"fmt"
	"math"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/server-selfish/backend/internal/domain/schema"
	monitoring_infra "github.com/server-selfish/backend/internal/infra/monitoring"
)

type (
	PrometheusRepository interface {
		GetCPUUsage(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetMemoryUsage(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetNetworkTx(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetNetworkRx(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetIORead(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
		GetIOWrite(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error)
	}
	prometheusRepository struct {
		pi monitoring_infra.PrometheusInfra
	}
)

func NewPrometheusRepository(pi monitoring_infra.PrometheusInfra) PrometheusRepository {
	return prometheusRepository{
		pi: pi,
	}
}

// GetCPUUsage implements [PrometheusRepository].
func (p prometheusRepository) GetCPUUsage(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	q := fmt.Sprintf(`rate(container_cpu_usage_seconds_total{name="%s"}[5m]) * 100`, params.ContainerName)
	res, _, err := p.pi.QueryRange(
		ctx,
		q,
		v1.Range{
			Start: params.StartTime,
			End:   params.EndTime,
			Step:  params.Step,
		},
	)
	if err != nil {
		return nil, err
	}

	matrix, ok := res.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("expected model.Matrix, got %T", res)
	}

	if len(matrix) == 0 {
		return []schema.MetricsSample{}, nil
	}

	samples := make([]schema.MetricsSample, len(matrix[0].Values))
	for i, sp := range matrix[0].Values {
		samples[i] = schema.MetricsSample{
			Timestamp: float64(sp.Timestamp),
			Value:     float64(math.Round(float64(sp.Value)*1000) / 1000),
		}
	}
	return samples, nil
}

// GetIORead implements [PrometheusRepository].
func (p prometheusRepository) GetIORead(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	q := fmt.Sprintf(`rate(container_fs_reads_bytes_total{name="%s"}[5m])`, params.ContainerName)
	res, _, err := p.pi.QueryRange(
		ctx,
		q,
		v1.Range{
			Start: params.StartTime,
			End:   params.EndTime,
			Step:  params.Step,
		},
	)
	if err != nil {
		return nil, err
	}

	matrix, ok := res.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("expected model.Matrix, got %T", res)
	}

	if len(matrix) == 0 {
		return []schema.MetricsSample{}, nil
	}
	samples := make([]schema.MetricsSample, len(matrix[0].Values))
	for i, sp := range matrix[0].Values {
		mbValue := float64(sp.Value) / (1024 * 1024)
		rounded := math.Round(mbValue*1000) / 1000

		samples[i] = schema.MetricsSample{
			Timestamp: float64(sp.Timestamp),
			Value:     rounded,
		}
	}
	return samples, nil
}

// GetIOWrite implements [PrometheusRepository].
func (p prometheusRepository) GetIOWrite(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	q := fmt.Sprintf(`rate(container_fs_writes_bytes_total{name="%s"}[5m])`, params.ContainerName)
	res, _, err := p.pi.QueryRange(
		ctx,
		q,
		v1.Range{
			Start: params.StartTime,
			End:   params.EndTime,
			Step:  params.Step,
		},
	)
	if err != nil {
		return nil, err
	}

	matrix, ok := res.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("expected model.Matrix, got %T", res)
	}

	if len(matrix) == 0 {
		return []schema.MetricsSample{}, nil
	}

	samples := make([]schema.MetricsSample, len(matrix[0].Values))
	for i, sp := range matrix[0].Values {
		mbValue := float64(sp.Value) / (1024 * 1024)
		rounded := math.Round(mbValue*1000) / 1000

		samples[i] = schema.MetricsSample{
			Timestamp: float64(sp.Timestamp),
			Value:     rounded,
		}
	}
	return samples, nil
}

// GetMemoryUsage implements [PrometheusRepository].
func (p prometheusRepository) GetMemoryUsage(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	q := fmt.Sprintf(`container_memory_usage_bytes{name="%s"}`, params.ContainerName)
	res, _, err := p.pi.QueryRange(
		ctx,
		q,
		v1.Range{
			Start: params.StartTime,
			End:   params.EndTime,
			Step:  params.Step,
		},
	)
	if err != nil {
		return nil, err
	}

	matrix, ok := res.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("expected model.Matrix, got %T", res)
	}

	if len(matrix) == 0 {
		return []schema.MetricsSample{}, nil
	}

	samples := make([]schema.MetricsSample, len(matrix[0].Values))
	for i, sp := range matrix[0].Values {
		mbValue := float64(sp.Value) / (1024 * 1024)
		rounded := math.Round(mbValue*1000) / 1000

		samples[i] = schema.MetricsSample{
			Timestamp: float64(sp.Timestamp),
			Value:     rounded,
		}
	}
	return samples, nil
}

// GetNetworkRx implements [PrometheusRepository].
func (p prometheusRepository) GetNetworkRx(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	q := fmt.Sprintf(`rate(container_network_receive_bytes_total{name="%s"}[5m])`, params.ContainerName)
	res, _, err := p.pi.QueryRange(
		ctx,
		q,
		v1.Range{
			Start: params.StartTime,
			End:   params.EndTime,
			Step:  params.Step,
		},
	)
	if err != nil {
		return nil, err
	}

	matrix, ok := res.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("expected model.Matrix, got %T", res)
	}

	if len(matrix) == 0 {
		return []schema.MetricsSample{}, nil
	}

	samples := make([]schema.MetricsSample, len(matrix[0].Values))
	for i, sp := range matrix[0].Values {
		mbValue := float64(sp.Value) / (1024 * 1024)
		rounded := math.Round(mbValue*1000) / 1000

		samples[i] = schema.MetricsSample{
			Timestamp: float64(sp.Timestamp),
			Value:     rounded,
		}
	}
	return samples, nil
}

// GetNetworkTx implements [PrometheusRepository].
func (p prometheusRepository) GetNetworkTx(ctx context.Context, params schema.GetQueryRangePrometheusRepositoryParams) ([]schema.MetricsSample, error) {
	q := fmt.Sprintf(`rate(container_network_transmit_bytes_total{name="%s"}[5m])`, params.ContainerName)
	res, _, err := p.pi.QueryRange(
		ctx,
		q,
		v1.Range{
			Start: params.StartTime,
			End:   params.EndTime,
			Step:  params.Step,
		},
	)
	if err != nil {
		return nil, err
	}

	matrix, ok := res.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("expected model.Matrix, got %T", res)
	}

	if len(matrix) == 0 {
		return []schema.MetricsSample{}, nil
	}

	samples := make([]schema.MetricsSample, len(matrix[0].Values))
	for i, sp := range matrix[0].Values {
		mbValue := float64(sp.Value) / (1024 * 1024)
		rounded := math.Round(mbValue*1000) / 1000

		samples[i] = schema.MetricsSample{
			Timestamp: float64(sp.Timestamp),
			Value:     rounded,
		}
	}
	return samples, nil
}
