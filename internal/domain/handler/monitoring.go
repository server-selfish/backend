package handler

import (
	"errors"
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
// GetCPUUsage  godoc
// @Summary     get cpu usage by container name
// @Tags        monitoring
// @Produce     json
// @Param       start_time query string true "start time (RFC3339)"
// @Param       end_time query string true "end time (RFC3339)"
// @Param       container_name query string true "container name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.MetricsReturn} "JSON response"
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /monitoring/cpu [get]
func (m monitoringHandler) GetCPUUsage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		m.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}
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
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		m.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	metrics, err := m.ps.GetCPUUsage(ctx, ui, schema.GetQueryRangePrometheusRepositoryParams{
		ContainerName: req.ContainerName,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		m.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "success", schema.MetricsReturn{Metrics: metrics})
}

// GetIORead implements [MonitoringHandler].
// GetIORead  godoc
// @Summary     get disk read IO by container name
// @Tags        monitoring
// @Produce     json
// @Param       start_time query string true "start time (RFC3339)"
// @Param       end_time query string true "end time (RFC3339)"
// @Param       container_name query string true "container name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.MetricsReturn} "JSON response"
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /monitoring/ioread [get]
func (m monitoringHandler) GetIORead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		m.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

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
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		m.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	metrics, err := m.ps.GetIORead(ctx, ui, schema.GetQueryRangePrometheusRepositoryParams{
		ContainerName: req.ContainerName,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		m.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "success", schema.MetricsReturn{Metrics: metrics})
}

// GetIOWrite implements [MonitoringHandler].
// GetIOWrite  godoc
// @Summary     get disk write IO by container name
// @Tags        monitoring
// @Produce     json
// @Param       start_time query string true "start time (RFC3339)"
// @Param       end_time query string true "end time (RFC3339)"
// @Param       container_name query string true "container name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.MetricsReturn} "JSON response"
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /monitoring/iowrite [get]
func (m monitoringHandler) GetIOWrite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		m.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

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
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		m.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	metrics, err := m.ps.GetIOWrite(ctx, ui, schema.GetQueryRangePrometheusRepositoryParams{
		ContainerName: req.ContainerName,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		m.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "success", schema.MetricsReturn{Metrics: metrics})
}

// GetMemoryUsage implements [MonitoringHandler].
// GetMemoryUsage  godoc
// @Summary     get memory usage by container name
// @Tags        monitoring
// @Produce     json
// @Param       start_time query string true "start time (RFC3339)"
// @Param       end_time query string true "end time (RFC3339)"
// @Param       container_name query string true "container name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.MetricsReturn} "JSON response"
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /monitoring/memory [get]
func (m monitoringHandler) GetMemoryUsage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		m.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

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
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		m.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	metrics, err := m.ps.GetMemoryUsage(ctx, ui, schema.GetQueryRangePrometheusRepositoryParams{
		ContainerName: req.ContainerName,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		m.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "success", schema.MetricsReturn{Metrics: metrics})
}

// GetNetworkRx implements [MonitoringHandler].
// GetNetworkRx  godoc
// @Summary     get network receive bytes by container name
// @Tags        monitoring
// @Produce     json
// @Param       start_time query string true "start time (RFC3339)"
// @Param       end_time query string true "end time (RFC3339)"
// @Param       container_name query string true "container name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.MetricsReturn} "JSON response"
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /monitoring/networkrx [get]
func (m monitoringHandler) GetNetworkRx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		m.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

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
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		m.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	metrics, err := m.ps.GetNetworkRx(ctx, ui, schema.GetQueryRangePrometheusRepositoryParams{
		ContainerName: req.ContainerName,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		m.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "success", schema.MetricsReturn{Metrics: metrics})
}

// GetNetworkTx implements [MonitoringHandler].
// GetNetworkTx  godoc
// @Summary     get network transmit bytes by container name
// @Tags        monitoring
// @Produce     json
// @Param       start_time query string true "start time (RFC3339)"
// @Param       end_time query string true "end time (RFC3339)"
// @Param       container_name query string true "container name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.MetricsReturn} "JSON response"
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /monitoring/networktx [get]
func (m monitoringHandler) GetNetworkTx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		m.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

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
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		m.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	metrics, err := m.ps.GetNetworkTx(ctx, ui, schema.GetQueryRangePrometheusRepositoryParams{
		ContainerName: req.ContainerName,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	})
	if err != nil {
		m.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "success", schema.MetricsReturn{Metrics: metrics})
}
