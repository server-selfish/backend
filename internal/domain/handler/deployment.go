package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

type (
	DeploymentHandler interface {
		GetDeploymentsByProjectId(w http.ResponseWriter, r *http.Request)
		GetDeploymentByDeploymentId(w http.ResponseWriter, r *http.Request)
		GetActiveDeploymentByDeploymentName(w http.ResponseWriter, r *http.Request)
		GetHistoryDeploymentByDeploymentName(w http.ResponseWriter, r *http.Request)
		GetTechstackName(w http.ResponseWriter, r *http.Request)
		GetTechstackVersionByName(w http.ResponseWriter, r *http.Request)
		GetDeploymentSettings(w http.ResponseWriter, r *http.Request)
		CreateNewDeployment(w http.ResponseWriter, r *http.Request)
		UpdateDeploymentVersionToLatest(w http.ResponseWriter, r *http.Request)
		UpdateDeployment(w http.ResponseWriter, r *http.Request)
		DeleteDeploymentByDeploymentId(w http.ResponseWriter, r *http.Request)
		DeleteDeploymentByDeploymentName(w http.ResponseWriter, r *http.Request)
	}
	deploymentHandler struct {
		ds     service.DeploymentService
		logger *zerolog.Logger
	}
)

func NewDeploymentHandler(ds service.DeploymentService, logger zerolog.Logger) DeploymentHandler {
	return &deploymentHandler{
		ds:     ds,
		logger: &logger,
	}
}

// DeleteDeploymentByDeploymentName implements [DeploymentHandler].
func (d *deploymentHandler) DeleteDeploymentByDeploymentName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	projectName := r.URL.Query().Get("project_name")
	deploymentName := r.URL.Query().Get("deployment_name")
	if projectName == "" {
		d.logger.Error().Msg(defined_error.ErrMissingProjectNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingProjectNameInParams)
		return
	}
	if deploymentName == "" {
		d.logger.Error().Msg(defined_error.ErrMissingDeploymentNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingDeploymentNameInParams)
		return
	}

	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	if err := d.ds.DeleteDeploymentByDeploymentName(ctx, ui, projectName, deploymentName); err != nil {
		d.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "deployment deleted", nil)
}

// UpdateDeploymentData implements [DeploymentHandler].
func (d *deploymentHandler) UpdateDeployment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}
	req, sc, err, ok := pkg.DecodeAndValidateBody[schema.UpdateDeploymentParams](w, r, d.logger)
	if !ok {
		pkg.ReturnError(w, sc, err)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	if err := d.ds.UpdateDeployment(ctx, ui, req); err != nil {
		d.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "settings updated", nil)
}

// GetDeploymentSetting implements [DeploymentHandler].
func (d *deploymentHandler) GetDeploymentSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	projectName := r.URL.Query().Get("project_name")
	deploymentName := r.URL.Query().Get("deployment_name")
	if projectName == "" {
		d.logger.Error().Msg(defined_error.ErrMissingProjectNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingProjectNameInParams)
		return
	}
	if deploymentName == "" {
		d.logger.Error().Msg(defined_error.ErrMissingDeploymentNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingDeploymentNameInParams)
		return
	}

	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	s, err := d.ds.GetDeploymentSettings(ctx, ui, projectName, deploymentName)
	if err != nil {
		d.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrActiveDeploymentNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch setting success", s)
}

// UpdateDeploymentVersionToLatest implements [DeploymentHandler].
func (d *deploymentHandler) UpdateDeploymentVersionToLatest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	req, sc, err, ok := pkg.DecodeAndValidateBody[schema.UpdateDeploymentHistoryToLatestParams](w, r, d.logger)
	if !ok {
		pkg.ReturnError(w, sc, err)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	if err := d.ds.UpdateDeploymentVersionToLatest(ctx, ui, req.ProjectName, req.DeploymentName); err != nil {
		d.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "new version deployed", nil)
}

// GetTechstackName implements [DeploymentHandler].
func (d *deploymentHandler) GetTechstackName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tl, err := d.ds.GetTechstackName(ctx)
	if err != nil {
		d.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch techstack sucess", tl)
}

// GetTechstackVersionByName implements [DeploymentHandler].
func (d *deploymentHandler) GetTechstackVersionByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tn := chi.URLParam(r, "techstack_name")
	if tn == "" {
		d.logger.Error().Msg(defined_error.ErrMissingTechstackNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingTechstackNameInParams)
		return
	}
	vl, err := d.ds.GetTechstackVersionByName(ctx, tn)
	if err != nil {
		d.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch version success", vl)
}

// CreateNewDeploymentVersion implements [DeploymentHandler].
func (d *deploymentHandler) CreateNewDeployment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}
	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	req, sc, err, ok := pkg.DecodeAndValidateBody[schema.CreateDeploymentHistoryParams](w, r, d.logger)
	if !ok {
		pkg.ReturnError(w, sc, err)
		return
	}

	ii, err := strconv.ParseInt(req.InstallationID, 10, 64)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringIntTypeCasting.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInvalidInstallationId)
		return
	}

	if err := d.ds.CreateNewDeployment(ctx, ui, ii, req); err != nil {
		d.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "deployment deployed", nil)
}

// DeleteDeploymentByDeploymentId implements [DeploymentHandler].
func (d *deploymentHandler) DeleteDeploymentByDeploymentId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		d.logger.Error().Msg(defined_error.ErrMissingIdInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingIdInParams)
		return
	}

	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}

	id, err := pkg.StringToPgUUID(idStr)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	if err := d.ds.DeleteDeploymentByDeploymentId(ctx, ui, id); err != nil {
		d.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "deployment deleted", nil)
}

// GetHistoryDeploymentByDeploymentId implements [DeploymentHandler].
func (d *deploymentHandler) GetHistoryDeploymentByDeploymentName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	projectName := r.URL.Query().Get("project_name")
	deploymentName := r.URL.Query().Get("deployment_name")
	if projectName == "" {
		d.logger.Error().Msg(defined_error.ErrMissingProjectNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingProjectNameInParams)
		return
	}
	if deploymentName == "" {
		d.logger.Error().Msg(defined_error.ErrMissingDeploymentNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingDeploymentNameInParams)
		return
	}

	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}

	deployments, err := d.ds.GetHistoryDeploymentByDeploymentName(ctx, ui, projectName, deploymentName)
	if err != nil {
		d.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", deployments)
}

// GetActiveDeploymenByDeploymentName implements [DeploymentHandler].
func (d *deploymentHandler) GetActiveDeploymentByDeploymentName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	projectName := r.URL.Query().Get("project_name")
	deploymentName := r.URL.Query().Get("deployment_name")
	if projectName == "" {
		d.logger.Error().Msg(defined_error.ErrMissingProjectNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingProjectNameInParams)
		return
	}
	if deploymentName == "" {
		d.logger.Error().Msg(defined_error.ErrMissingDeploymentNameInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingDeploymentNameInParams)
		return
	}

	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}

	deployment, err := d.ds.GetActiveDeploymentByDeploymentName(ctx, ui, projectName, deploymentName)
	if err != nil {
		d.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrActiveDeploymentNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", deployment)
}

// GetDeploymentByDeploymentId implements [DeploymentHandler].
func (d *deploymentHandler) GetDeploymentByDeploymentId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		d.logger.Error().Msg(defined_error.ErrMissingIdInParams.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissingIdInParams)
		return
	}

	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}

	id, err := pkg.StringToPgUUID(idStr)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	deployment, err := d.ds.GetDeploymentByDeploymentId(ctx, ui, id)
	if err != nil {
		d.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrDeploymentNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", deployment)
}

// GetDeploymentsByProjectId implements [DeploymentHandler].
func (d *deploymentHandler) GetDeploymentsByProjectId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	q := r.URL.Query()
	projectId := q.Get("project_id")
	if projectId == "" {
		d.logger.Error().Msg(defined_error.ErrMissinProjectId.Error())
		pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrMissinProjectId)
		return
	}

	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}

	id, err := pkg.StringToPgUUID(projectId)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	deployments, err := d.ds.GetDeploymentsByProjectId(ctx, ui, id)
	if err != nil {
		d.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", deployments)
}
