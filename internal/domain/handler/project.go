package handler

import (
	"net/http"

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
		ps service.ProjectService
	}
)

func NewProjectHandler(ps service.ProjectService) ProjectHandler {
	return &projectHandler{
		ps: ps,
	}
}

// DeleteProjectById implements [ProjectHandler].
func (p *projectHandler) DeleteProjectById(w http.ResponseWriter, r *http.Request) {
	req, ok := pkg.DecodeAndValidate[schema.DeleteProjectParams](w, r)
	if !ok {
		return
	}
	ctx := r.Context()
	id, err := pkg.StringToPgUUID(req.ID)
	if err != nil {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	if err := p.ps.DeleteProjectById(ctx, id); err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "project deleted", nil)
}

// GetAllProjects implements [ProjectHandler].
func (p *projectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projects, err := p.ps.GetAllProjects(ctx)
	if err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", projects)
}

// GetProjectById implements [ProjectHandler].
func (p *projectHandler) GetProjectById(w http.ResponseWriter, r *http.Request) {
	req, ok := pkg.DecodeAndValidate[schema.GetProjectParams](w, r)
	if !ok {
		return
	}
	ctx := r.Context()
	id, err := pkg.StringToPgUUID(req.ID)
	if err != nil {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	project, err := p.ps.GetProjectById(ctx, id)
	if err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusCreated, "fetch success", project)
}

// UpdateProjectById implements [ProjectHandler].
func (p *projectHandler) UpdateProjectById(w http.ResponseWriter, r *http.Request) {
	req, ok := pkg.DecodeAndValidate[schema.UpdateProjectParams](w, r)
	if !ok {
		return
	}
	ctx := r.Context()
	_, err := pkg.StringToPgUUID(req.ID)
	if err != nil {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	if err := p.ps.UpdateProjectById(ctx, &req); err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "project updated", nil)
}

// CreateProject implements [ProjectHandler].
func (p *projectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	req, ok := pkg.DecodeAndValidate[schema.CreateProjectParams](w, r)
	if !ok {
		return
	}
	ctx := r.Context()
	if err := p.ps.CreateProject(ctx, &req); err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusCreated, "Project Created", nil)
}
