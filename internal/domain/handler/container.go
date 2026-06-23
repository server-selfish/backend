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

	events, errs, qErr := c.cs.StreamLogs(c.appCtx, ui, name)
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
	for {
		select {
		case <-c.appCtx.Done():
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
