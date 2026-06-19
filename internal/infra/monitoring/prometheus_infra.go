package monitoring_infra

import (
	"context"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type (
	PrometheusInfra interface {
		v1.API
	}
	prometheusInfra struct {
		pc v1.API
	}
)

func NewPrometheusInfra(pc v1.API) PrometheusInfra {
	return prometheusInfra{
		pc: pc,
	}
}

// AlertManagers implements [PrometheusInfra].
func (p prometheusInfra) AlertManagers(ctx context.Context) (v1.AlertManagersResult, error) {
	panic("unimplemented")
}

// Alerts implements [PrometheusInfra].
func (p prometheusInfra) Alerts(ctx context.Context) (v1.AlertsResult, error) {
	panic("unimplemented")
}

// Buildinfo implements [PrometheusInfra].
func (p prometheusInfra) Buildinfo(ctx context.Context) (v1.BuildinfoResult, error) {
	panic("unimplemented")
}

// CleanTombstones implements [PrometheusInfra].
func (p prometheusInfra) CleanTombstones(ctx context.Context) error {
	panic("unimplemented")
}

// Config implements [PrometheusInfra].
func (p prometheusInfra) Config(ctx context.Context) (v1.ConfigResult, error) {
	panic("unimplemented")
}

// DeleteSeries implements [PrometheusInfra].
func (p prometheusInfra) DeleteSeries(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) error {
	panic("unimplemented")
}

// Flags implements [PrometheusInfra].
func (p prometheusInfra) Flags(ctx context.Context) (v1.FlagsResult, error) {
	panic("unimplemented")
}

// LabelNames implements [PrometheusInfra].
func (p prometheusInfra) LabelNames(ctx context.Context, matches []string, startTime time.Time, endTime time.Time, opts ...v1.Option) ([]string, v1.Warnings, error) {
	panic("unimplemented")
}

// LabelValues implements [PrometheusInfra].
func (p prometheusInfra) LabelValues(ctx context.Context, label string, matches []string, startTime time.Time, endTime time.Time, opts ...v1.Option) (model.LabelValues, v1.Warnings, error) {
	panic("unimplemented")
}

// Metadata implements [PrometheusInfra].
func (p prometheusInfra) Metadata(ctx context.Context, metric string, limit string) (map[string][]v1.Metadata, error) {
	panic("unimplemented")
}

// Query implements [PrometheusInfra].
func (p prometheusInfra) Query(ctx context.Context, query string, ts time.Time, opts ...v1.Option) (model.Value, v1.Warnings, error) {
	return p.pc.Query(ctx, query, ts, opts...)
}

// QueryExemplars implements [PrometheusInfra].
func (p prometheusInfra) QueryExemplars(ctx context.Context, query string, startTime time.Time, endTime time.Time) ([]v1.ExemplarQueryResult, error) {
	return p.pc.QueryExemplars(ctx, query, startTime, endTime)
}

// QueryRange implements [PrometheusInfra].
func (p prometheusInfra) QueryRange(ctx context.Context, query string, r v1.Range, opts ...v1.Option) (model.Value, v1.Warnings, error) {
	return p.pc.QueryRange(ctx, query, r, opts...)
}

// Rules implements [PrometheusInfra].
func (p prometheusInfra) Rules(ctx context.Context) (v1.RulesResult, error) {
	panic("unimplemented")
}

// Runtimeinfo implements [PrometheusInfra].
func (p prometheusInfra) Runtimeinfo(ctx context.Context) (v1.RuntimeinfoResult, error) {
	panic("unimplemented")
}

// Series implements [PrometheusInfra].
func (p prometheusInfra) Series(ctx context.Context, matches []string, startTime time.Time, endTime time.Time, opts ...v1.Option) ([]model.LabelSet, v1.Warnings, error) {
	panic("unimplemented")
}

// Snapshot implements [PrometheusInfra].
func (p prometheusInfra) Snapshot(ctx context.Context, skipHead bool) (v1.SnapshotResult, error) {
	panic("unimplemented")
}

// TSDB implements [PrometheusInfra].
func (p prometheusInfra) TSDB(ctx context.Context, opts ...v1.Option) (v1.TSDBResult, error) {
	panic("unimplemented")
}

// Targets implements [PrometheusInfra].
func (p prometheusInfra) Targets(ctx context.Context) (v1.TargetsResult, error) {
	panic("unimplemented")
}

// TargetsMetadata implements [PrometheusInfra].
func (p prometheusInfra) TargetsMetadata(ctx context.Context, matchTarget string, metric string, limit string) ([]v1.MetricMetadata, error) {
	panic("unimplemented")
}

// WalReplay implements [PrometheusInfra].
func (p prometheusInfra) WalReplay(ctx context.Context) (v1.WalReplayStatus, error) {
	panic("unimplemented")
}
