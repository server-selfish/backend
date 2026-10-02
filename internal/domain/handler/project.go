package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

type (
	ProjectHandler interface {
		CreateProject(w http.ResponseWriter, r *http.Request)
		GetAllProjects(w http.ResponseWriter, r *http.Request)
		// GetProjectById(w http.ResponseWriter, r *http.Request)
		GetProjectByName(w http.ResponseWriter, r *http.Request)
		GetProjectByNameDetail(w http.ResponseWriter, r *http.Request)
		// UpdateProjectById(w http.ResponseWriter, r *http.Request)
		DeleteProjectById(w http.ResponseWriter, r *http.Request)
	}
	projectHandler struct {
		ps     service.ProjectService
		logger *zerolog.Logger
	}
)

func NewProjectHandler(ps service.ProjectService, logger zerolog.Logger) ProjectHandler {
	return &projectHandler{
		ps:     ps,
		logger: &logger,
	}
}

// GetProjectByNameDetail implements [ProjectHandler].
// GetProjectByNameDetail godoc
// @Summary     Get project detail with its deployment, by its name
// @Tags        projects
// @Produce     json
// @Param       name path string true "Project name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.GetProjectAllDetail}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /project/detail/{name} [get]
func (p *projectHandler) GetProjectByNameDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		p.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}
	name := chi.URLParam(r, "name")
	if name == "" {
		p.logger.Error().Msg(defined_error.ErrMissingNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingNameInParams)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		p.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	project, err := p.ps.GetProjectByNameDetail(ctx, name, ui)
	if err != nil {
		p.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrProjectNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", project)
}

// GetProjectByName implements [ProjectHandler].
// GetProjectByName godoc
// @Summary     Get project by its name
// @Tags        projects
// @Produce     json
// @Param       name path string true "Project name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.GetProjectDetail}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /project/{name} [get]
func (p *projectHandler) GetProjectByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		p.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}
	name := chi.URLParam(r, "name")
	if name == "" {
		p.logger.Error().Msg(defined_error.ErrMissingNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingNameInParams)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		p.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	project, err := p.ps.GetProjectByName(ctx, name, ui)
	if err != nil {
		p.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrProjectNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", project)
}

// DeleteProjectById implements [ProjectHandler].
// DeleteProjectById godoc
// @Summary     Delete project by its id
// @Tags        projects
// @Produce     json
// @Param       id path string true "Project id"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{message="project deleted"}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /project/{id} [delete]
func (p *projectHandler) DeleteProjectById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		p.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		p.logger.Error().Msg(defined_error.ErrMissingIdInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingIdInParams)
		return
	}
	id, err := pkg.StringToPgUUID(idStr)
	if err != nil {
		p.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}

	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		p.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	if err := p.ps.DeleteProjectById(ctx, id, ui); err != nil {
		p.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "project deleted", nil)
}

// GetAllProjects implements [ProjectHandler].
// GetAllProjects godoc
// @Summary     Get all projects
// @Tags        projects
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=[]schema.GetProjectsData}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /project [get]
func (p *projectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		p.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}
	id, err := pkg.StringToPgUUID(userID)
	if err != nil {
		p.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	projects, err := p.ps.GetAllProjects(ctx, id)
	if err != nil {
		p.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", projects)
}

// CreateProject implements [ProjectHandler].
// CreateProject godoc
// @Summary     Create new project
// @Tags        projects
// @Accept 			json
// @Produce     json
// @Param       request body schema.CreateProjectParams true "Create project request"
// @Security    BearerAuth
// @Success     201 {object} pkg.Response{message="Project Created"}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     409 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /project [post]
func (p *projectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		p.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		p.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	req, sc, err, ok := pkg.DecodeAndValidateBody[schema.CreateProjectParams](w, r, p.logger)
	if !ok {
		pkg.ReturnError(w, sc, err)
		return
	}
	if err := p.ps.CreateProject(ctx, ui, &req); err != nil {
		p.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrProjectNameUniqueConstraint):
			pkg.ReturnError(w, http.StatusConflict, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusCreated, "Project Created", nil)
}
