package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
)

type (
	ProjectHandler interface {
		CreateProject(w http.ResponseWriter, r *http.Request)
		GetAllProjects(w http.ResponseWriter, r *http.Request)
		GetProjectById(w http.ResponseWriter, r *http.Request)
		UpdateProjectById(w http.ResponseWriter, r *http.Request)
		DeleteProjectById(w http.ResponseWriter, r *http.Request)
	}
	projectHandler struct {
		ps     service.ProjectService
		logger zerolog.Logger
	}
)

func NewProjectHandler(ps service.ProjectService, logger zerolog.Logger) ProjectHandler {
	return &projectHandler{
		ps:     ps,
		logger: logger,
	}
}

// DeleteProjectById implements [ProjectHandler].
func (p *projectHandler) DeleteProjectById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "unauthorized",
		})
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	id, err := pkg.StringToPgUUID(idStr)
	if err != nil {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}

	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	if err := p.ps.DeleteProjectById(ctx, id, ui); err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "project deleted", nil)
}

// GetAllProjects implements [ProjectHandler].
func (p *projectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		pkg.WriteJSON(w, http.StatusUnauthorized, schema.AuthErrorResponse{
			Message: "unauthorized",
		})
		return
	}
	id, err := pkg.StringToPgUUID(userID)
	if err != nil {
		p.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	projects, err := p.ps.GetAllProjects(ctx, id)
	if err != nil {
		p.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", projects)
}

// GetProjectById implements [ProjectHandler].
func (p *projectHandler) GetProjectById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	ctx := r.Context()
	id, err := pkg.StringToPgUUID(idStr)
	if err != nil {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
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
	project, err := p.ps.GetProjectById(ctx, id, ui)
	if err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusCreated, "fetch success", project)
}

// UpdateProjectById implements [ProjectHandler].
func (p *projectHandler) UpdateProjectById(w http.ResponseWriter, r *http.Request) {
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
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	id, err := pkg.StringToPgUUID(idStr)
	if err != nil {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	req, ok := pkg.DecodeAndValidateBody[schema.UpdateProjectParams](w, r)
	if !ok {
		return
	}
	if err := p.ps.UpdateProjectById(ctx, id, ui, &req); err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "project updated", nil)
}

// CreateProject implements [ProjectHandler].
func (p *projectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
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
		p.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	req, ok := pkg.DecodeAndValidateBody[schema.CreateProjectParams](w, r)
	if !ok {
		return
	}
	if err := p.ps.CreateProject(ctx, ui, &req); err != nil {
		p.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusCreated, "Project Created", nil)
}
