package handler

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

type (
	MonitoringHandler interface {
		GetCPUUsage(w http.ResponseWriter, r *http.Request)
		GetIORead(w http.ResponseWriter, r *http.Request)
		GetIOWrite(w http.ResponseWriter, r *http.Request)
		GetMemoryUsage(w http.ResponseWriter, r *http.Request)
		GetNetworkRx(w http.ResponseWriter, r *http.Request)
		GetNetworkTx(w http.ResponseWriter, r *http.Request)
	}
	monitoringHandler struct {
		ps     service.PrometheusService
		logger zerolog.Logger
	}
)

func NewMonitoringHandler(ps service.PrometheusService, logger zerolog.Logger) MonitoringHandler {
	return monitoringHandler{
		ps:     ps,
		logger: logger,
	}
}

// GetCPUUsage implements [MonitoringHandler].
func (m monitoringHandler) GetCPUUsage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req schema.MetricsRequest

	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidStartTimeFormat)
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidEndTimeFormat)
		return
	}

	req.ContainerName = r.URL.Query().Get("container_name")
	req.StartTime = startTime
	req.EndTime = endTime

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		m.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metrics, err := m.ps.GetCPUUsage(ctx, schema.GetQueryRangePrometheusRepositoryParams{
		ContainerName: req.ContainerName,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "success", schema.MetricsReturn{Metrics: metrics})
}

// GetIORead implements [MonitoringHandler].
func (m monitoringHandler) GetIORead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req schema.MetricsRequest

	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidStartTimeFormat)
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidEndTimeFormat)
		return
	}

	req.ContainerName = r.URL.Query().Get("container_name")
	req.StartTime = startTime
	req.EndTime = endTime

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		m.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metrics, err := m.ps.GetIORead(ctx, schema.GetQueryRangePrometheusRepositoryParams{
		ContainerName: req.ContainerName,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "success", schema.MetricsReturn{Metrics: metrics})
}

// GetIOWrite implements [MonitoringHandler].
func (m monitoringHandler) GetIOWrite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req schema.MetricsRequest

	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidStartTimeFormat)
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidEndTimeFormat)
		return
	}

	req.ContainerName = r.URL.Query().Get("container_name")
	req.StartTime = startTime
	req.EndTime = endTime

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		m.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metrics, err := m.ps.GetIOWrite(ctx, schema.GetQueryRangePrometheusRepositoryParams{
		ContainerName: req.ContainerName,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "success", schema.MetricsReturn{Metrics: metrics})
}

// GetMemoryUsage implements [MonitoringHandler].
func (m monitoringHandler) GetMemoryUsage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req schema.MetricsRequest

	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidStartTimeFormat)
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidEndTimeFormat)
		return
	}

	req.ContainerName = r.URL.Query().Get("container_name")
	req.StartTime = startTime
	req.EndTime = endTime

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		m.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metrics, err := m.ps.GetMemoryUsage(ctx, schema.GetQueryRangePrometheusRepositoryParams{
		ContainerName: req.ContainerName,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "success", schema.MetricsReturn{Metrics: metrics})
}

// GetNetworkRx implements [MonitoringHandler].
func (m monitoringHandler) GetNetworkRx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req schema.MetricsRequest

	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidStartTimeFormat)
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidEndTimeFormat)
		return
	}

	req.ContainerName = r.URL.Query().Get("container_name")
	req.StartTime = startTime
	req.EndTime = endTime

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		m.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metrics, err := m.ps.GetNetworkRx(ctx, schema.GetQueryRangePrometheusRepositoryParams{
		ContainerName: req.ContainerName,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "success", schema.MetricsReturn{Metrics: metrics})
}

// GetNetworkTx implements [MonitoringHandler].
func (m monitoringHandler) GetNetworkTx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req schema.MetricsRequest

	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidStartTimeFormat)
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidEndTimeFormat)
		return
	}

	req.ContainerName = r.URL.Query().Get("container_name")
	req.StartTime = startTime
	req.EndTime = endTime

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		m.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metrics, err := m.ps.GetNetworkTx(ctx, schema.GetQueryRangePrometheusRepositoryParams{
		ContainerName: req.ContainerName,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		m.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "success", schema.MetricsReturn{Metrics: metrics})
}
