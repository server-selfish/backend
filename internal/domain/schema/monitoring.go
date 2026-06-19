package schema

import (
	"time"
)

type (
	GetQueryRangePrometheusRepositoryParams struct {
		ContainerName string
		StartTime     time.Time
		EndTime       time.Time
		Step          time.Duration
	}
	MetricsRequest struct {
		ContainerName string    `json:"container_name" query:"container_name" validate:"required"`
		StartTime     time.Time `json:"startTime" query:"start_time" validate:"required"`
		EndTime       time.Time `json:"endTime" query:"end_time" validate:"required,gtfield=StartTime"`
	}
)

type (
	MetricsSample struct {
		Timestamp float64 `json:"timestamp"`
		Value     float64 `json:"value"`
	}
	MetricsReturn struct {
		Metrics []MetricsSample `json:"metrics"`
	}
)
