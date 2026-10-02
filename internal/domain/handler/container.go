package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
	"github.com/spf13/viper"
)

type (
	ContainerHandler interface {
		GetContainerStatus(w http.ResponseWriter, r *http.Request)
		PauseContainer(w http.ResponseWriter, r *http.Request)
		UnPauseContainer(w http.ResponseWriter, r *http.Request)
		StopContainer(w http.ResponseWriter, r *http.Request)
		StartContainer(w http.ResponseWriter, r *http.Request)
		RestartContainer(w http.ResponseWriter, r *http.Request)
		StreamLogs(w http.ResponseWriter, r *http.Request)
	}
	containerHandler struct {
		cs     service.ContainerService
		logger *zerolog.Logger
		appCtx context.Context
	}
)

func NewContainerHandler(cs service.ContainerService, logger zerolog.Logger, appCtx context.Context) ContainerHandler {
	return &containerHandler{
		cs:     cs,
		logger: &logger,
		appCtx: appCtx,
	}
}

// StreamLogs implements [ContainerHandler].
// StreamLogs godoc
// @Summary     Stream container logs
// @Description Streams container logs using Server-Sent Events (SSE).
// @Tags        containers
// @Produce     text/event-stream
// @Param       name path string true "Container name"
// @Security    BearerAuth
// @Success     200 {string} string "Server-Sent Events stream"
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /container/log/{name} [get]
func (c *containerHandler) StreamLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		c.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		c.logger.Error().Msg(defined_error.ErrMissingNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingNameInParams)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		c.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		c.logger.Error().Msg("error flush")
		pkg.ReturnError(w, http.StatusInternalServerError, errors.New("steraming unsupported"))
		return
	}

	// Bound the stream lifetime and propagate client disconnects to the
	// service so its goroutines and the docker log reader are released.
	// (Previously the app-wide context was passed, which outlives requests.)
	streamCtx, cancel := context.WithTimeout(r.Context(), resolveContainerLogSSETimeout())
	defer cancel()

	events, errs, qErr := c.cs.StreamLogs(streamCtx, ui, name)
	if qErr != nil {
		c.logger.Error().Msg(qErr.Error())
		switch {
		case errors.Is(qErr, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, qErr)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	// Prevent the nginx gateway from buffering the stream.
	w.Header().Set("X-Accel-Buffering", "no")
	for {
		select {
		case <-streamCtx.Done():
			return

		case err := <-errs:
			if err != nil {
				pkg.ReturnError(w, http.StatusInternalServerError, err)
			}
			return

		case event, ok := <-events:
			if !ok {
				return
			}

			b, _ := json.Marshal(event)

			if tcpConn, ok := w.(interface{ SetWriteDeadline(time.Time) }); ok {
				tcpConn.SetWriteDeadline(time.Now().Add(time.Second))
			}

			if _, err := fmt.Fprintf(w, "data: %s\n\n", b); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

// RestartContainer implements [ContainerHandler].
// RestartContainer godoc
// @Summary     		restart container based on its name
// @Tags        		containers
// @Produce     		json
// @Param       		name path string true "Container name"
// @Security    		BearerAuth
// @Success     		200 {object} pkg.Response{message="deployment restarted"}
// @Failure     		400 {object} pkg.Response{error=string}
// @Failure     		401 {object} pkg.Response{error=string}
// @Failure     		404 {object} pkg.Response{error=string}
// @Failure     		500 {object} pkg.Response{error=string}
// @Router      		/container/restart/{name} [post]
func (c *containerHandler) RestartContainer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		c.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		c.logger.Error().Msg(defined_error.ErrMissingNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingNameInParams)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		c.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	if err := c.cs.RestartContainer(ctx, name, ui); err != nil {
		c.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "deployment restarted", nil)
}

// PauseContainer implements [ContainerHandler].
// PauseContainer godoc
// @Summary     		pause container based on its name
// @Tags        		containers
// @Produce     		json
// @Param       		name path string true "Container name"
// @Security    		BearerAuth
// @Success     		200 {object} pkg.Response{message="deployment paused"}
// @Failure     		400 {object} pkg.Response{error=string}
// @Failure     		401 {object} pkg.Response{error=string}
// @Failure     		404 {object} pkg.Response{error=string}
// @Failure     		500 {object} pkg.Response{error=string}
// @Router      		/container/pause/{name} [post]
func (c *containerHandler) PauseContainer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		c.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		c.logger.Error().Msg(defined_error.ErrMissingNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingNameInParams)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		c.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	if err := c.cs.PauseContainer(ctx, name, ui); err != nil {
		c.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "deployment paused", nil)
}

// UnPauseContainer implements [ContainerHandler].
// UnPauseContainer godoc
// @Summary     		unpause container based on its name
// @Tags        		containers
// @Produce     		json
// @Param       		name path string true "Container name"
// @Security    		BearerAuth
// @Success     		200 {object} pkg.Response{message="deployment unpaused"}
// @Failure     		400 {object} pkg.Response{error=string}
// @Failure     		401 {object} pkg.Response{error=string}
// @Failure     		404 {object} pkg.Response{error=string}
// @Failure     		500 {object} pkg.Response{error=string}
// @Router      		/container/unpause/{name} [post]
func (c *containerHandler) UnPauseContainer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		c.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		c.logger.Error().Msg(defined_error.ErrMissingNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingNameInParams)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		c.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	if err := c.cs.UnPauseContainer(ctx, name, ui); err != nil {
		c.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "deployment unpaused", nil)
}

// StartContainer implements [ContainerHandler].
// StartContainer godoc
// @Summary     		start container based on its name
// @Tags        		containers
// @Produce     		json
// @Param       		name path string true "Container name"
// @Security    		BearerAuth
// @Success     		200 {object} pkg.Response{message="deployment started"}
// @Failure     		400 {object} pkg.Response{error=string}
// @Failure     		401 {object} pkg.Response{error=string}
// @Failure     		404 {object} pkg.Response{error=string}
// @Failure     		500 {object} pkg.Response{error=string}
// @Router      		/container/start/{name} [post]
func (c *containerHandler) StartContainer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		c.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		c.logger.Error().Msg(defined_error.ErrMissingNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingNameInParams)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		c.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	if err := c.cs.StartContainer(ctx, name, ui); err != nil {
		c.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "deployment started", nil)
}

// StopContainer implements [ContainerHandler].
// StopContainer godoc
// @Summary     		stop container based on its name
// @Tags        		containers
// @Produce     		json
// @Param       		name path string true "Container name"
// @Security    		BearerAuth
// @Success     		200 {object} pkg.Response{message="deployment stopped"}
// @Failure     		400 {object} pkg.Response{error=string}
// @Failure     		401 {object} pkg.Response{error=string}
// @Failure     		404 {object} pkg.Response{error=string}
// @Failure     		500 {object} pkg.Response{error=string}
// @Router      		/container/stop/{name} [post]
func (c *containerHandler) StopContainer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		c.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		c.logger.Error().Msg(defined_error.ErrMissingNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingNameInParams)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		c.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	if err := c.cs.StopContainer(ctx, name, ui); err != nil {
		c.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "deployment stopped", nil)
}

// GetContainerStatus implements [ContainerHandler].
// GetContainerStatus godoc
// @Summary     		get container status based on its name
// @Tags        		containers
// @Produce     		json
// @Param       		name path string true "Container name"
// @Security    		BearerAuth
// @Success     		200 {object} pkg.Response{data=schema.ContainerStatusResponse}
// @Failure     		400 {object} pkg.Response{error=string}
// @Failure     		401 {object} pkg.Response{error=string}
// @Failure     		404 {object} pkg.Response{error=string}
// @Failure     		500 {object} pkg.Response{error=string}
// @Router      		/container/status/{name} [get]
func (c *containerHandler) GetContainerStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		c.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		c.logger.Error().Msg(defined_error.ErrMissingNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingNameInParams)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		c.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	status, err := c.cs.GetContainerStatus(ctx, name, ui)
	if err != nil {
		c.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrContainerNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch status success", status)
}

// resolveContainerLogSSETimeout reads sse.containerlog.timeout_seconds with
// a 1800s fallback for missing or non-positive values. Live tails run far
// longer than build-log replays, hence the larger default.
func resolveContainerLogSSETimeout() time.Duration {
	viper.SetDefault("sse.containerlog.timeout_seconds", 1800)
	if n := viper.GetInt("sse.containerlog.timeout_seconds"); n > 0 {
		return time.Duration(n) * time.Second
	}
	return 1800 * time.Second
}
