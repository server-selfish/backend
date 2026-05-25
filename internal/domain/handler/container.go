package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

type (
	ContainerHandler interface {
		GetContainerStatus(w http.ResponseWriter, r *http.Request)
	}
	containerHandler struct {
		cs     service.ContainerService
		logger *zerolog.Logger
	}
)

func NewContainerHandler(cs service.ContainerService, logger zerolog.Logger) ContainerHandler {
	return &containerHandler{
		cs:     cs,
		logger: &logger,
	}
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
