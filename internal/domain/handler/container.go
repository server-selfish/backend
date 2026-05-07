package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
)

type (
	ContainerHandler interface {
		GetContainerStatus(w http.ResponseWriter, r *http.Request)
	}
	containerHandler struct {
		cs service.ContainerService
	}
)

func NewContainerHandler(cs service.ContainerService) ContainerHandler {
	return &containerHandler{
		cs: cs,
	}
}

// GetContainerStatus implements [ContainerHandler].
func (c *containerHandler) GetContainerStatus(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "unauthorized",
		})
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	status, err := c.cs.GetContainerStatus(ctx, name, ui)
	if err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch status success", status)
}
