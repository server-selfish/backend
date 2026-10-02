package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/internal/domain/schema"
	service "github.com/server-selfish/backend/internal/domain/service/deployment"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
	"github.com/spf13/viper"
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
		GetBuildLog(w http.ResponseWriter, r *http.Request)
		IsBuildRunning(w http.ResponseWriter, r *http.Request)
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
// DeleteDeploymentByDeploymentName godoc
// @Summary     Delete deployment by its name
// @Tags        deployments
// @Produce     json
// @Param       project_name query string true "Project name"
// @Param       deployment_name query string true "deployment name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{message="deployment deleted"}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment [delete]
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

// UpdateDeployment implements [DeploymentHandler].
// UpdateDeployment godoc
// @Summary     update deployment by its name
// @Tags        deployments
// @Accept      json
// @Produce     json
// @Param       request body schema.UpdateDeploymentParams true "update deployment request"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{message="settings updated"}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     409 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment [patch]
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
		if errors.Is(err, defined_error.ErrBuildInProgress) {
			pkg.ReturnError(w, http.StatusConflict, err)
			return
		}
		pkg.ReturnError(w, http.StatusInternalServerError, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "settings updated", nil)
}

// GetDeploymentSettings implements [DeploymentHandler].
// GetDeploymentSettings godoc
// @Summary     get deployment setting by its name
// @Tags        deployments
// @Produce     json
// @Param       project_name query string true "Project name"
// @Param       deployment_name query string true "deployment name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.GetDeploymentSettings}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment/settings [get]
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
// UpdateDeploymentVersionToLatest godoc
// @Summary     Re deploy the latest version of configured branch.
// @Tags        deployments
// @Accept	json
// @Produce     json
// @Param       request body schema.UpdateDeploymentHistoryToLatestParams true "update deployment request"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{message="new version deployed"}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     409 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment/deploy [post]
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
		if errors.Is(err, defined_error.ErrBuildInProgress) {
			pkg.ReturnError(w, http.StatusConflict, err)
			return
		}
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "new version deployed", nil)
}

// GetTechstackName implements [DeploymentHandler].
// GetTechstackName godoc
// @Summary     get techstack name list available
// @Tags        deployments
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.GetTechstackList}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment/techstack [get]
func (d *deploymentHandler) GetTechstackName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}
	tl, err := d.ds.GetTechstackName(ctx)
	if err != nil {
		d.logger.Error().Msg(err.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch techstack sucess", tl)
}

// GetTechstackVersionByName implements [DeploymentHandler].
// GetTechstackVersionByName godoc
// @Summary     get list of techstack version by name.
// @Tags        deployments
// @Produce     json
// @Param	techstack_name path string true "Techstack name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=[]schema.GetTechstackVersion}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment/techstack/{techstack_name}/version [get]
func (d *deploymentHandler) GetTechstackVersionByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}
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

// CreateNewDeployment implements [DeploymentHandler].
// CreateNewDeployment godoc
// @Summary     create new deployment.
// @Tags        deployments
// @Accept	json
// @Produce     json
// @Param	request body schema.CreateDeploymentHistoryParams true "Create deployment history request"
// @Security    BearerAuth
// @Success     201 {object} pkg.Response{message="deployment deployed"}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     409 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment [post]
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
		if errors.Is(err, defined_error.ErrBuildInProgress) {
			pkg.ReturnError(w, http.StatusConflict, err)
			return
		}
		pkg.ReturnError(w, http.StatusInternalServerError, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusCreated, "deployment deployed", nil)
}

// DeleteDeploymentByDeploymentId implements [DeploymentHandler].
// DeleteDeploymentByDeploymentId godoc
// @Summary     Delete deployment by deployment id.
// @Tags        deployments
// @Produce     json
// @Param	id path string true "deployment id"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{message="deployment deleted"}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment/{id} [delete]
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

// GetHistoryDeploymentByDeploymentName implements [DeploymentHandler].
// GetHistoryDeploymentByDeploymentName godoc
// @Summary     get deployment history by deployment name.
// @Tags        deployments
// @Produce     json
// @Param				project_name query string true "project name"
// @Param				deployment_name query string true "deployment name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=[]schema.GetHistoryDeploymentHistory}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment/history [get]
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
// GetActiveDeploymenByDeploymentName godoc
// @Summary     get active deployment by deployment name.
// @Tags        deployments
// @Produce     json
// @Param				project_name query string true "project name"
// @Param				deployment_name query string true "deployment name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=[]schema.GetHistoryDeploymentHistory}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment/active [get]
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
// GetDeploymentByDeploymentId godoc
// @Summary     get deployment by deploymet id.
// @Tags        deployments
// @Produce     json
// @Param				id path string true "deployment id"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.GetSingleDeploymentData}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment/{id} [get]
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
// GetDeploymentsByProjectId godoc
// @Summary     get all deployment by project id.
// @Tags        deployments
// @Produce     json
// @Param				project_id query string true "project id"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=[]schema.GetDeploymentData}
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment [get]
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

// GetBuildLog implements [DeploymentHandler].
// GetBuildLog godoc
// @Summary     get build log [for debug only].
// @Description Returns build log lines as JSON, or as a Server-Sent Events stream when the client sends Accept: text/event-stream.
// @Tags        deployments
// @Produce     json
// @Param       project_name query string true "project name"
// @Param       deployment_name query string true "deployment name"
// @Param       attempt_id query string false "deployment attempt id"
// @Param       limit query int false "maximum number of log lines" default(1000)
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=[]schema.BuildLogLine} "JSON response, or SSE event stream when Accept is text/event-stream"
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     404 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment/build-log [get]
func (d *deploymentHandler) GetBuildLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	q := r.URL.Query()
	projectName := q.Get("project_name")
	deploymentName := q.Get("deployment_name")
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

	limit := 1000
	if raw := q.Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			d.logger.Error().Msg("invalid limit query param")
			pkg.ReturnError(w, http.StatusBadRequest, defined_error.ErrInternalServerError)
			return
		}
		limit = n
	}

	ui, err := pkg.StringToPgUUID(userID)
	if err != nil {
		d.logger.Error().Err(err).Msg(defined_error.ErrStringUUIDTypeCasting.Error())
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}
	lines, err := d.ds.GetBuildLog(ctx, ui, projectName, deploymentName, q.Get("attempt_id"), limit)
	if err != nil {
		d.logger.Error().Msg(err.Error())
		switch {
		case errors.Is(err, defined_error.ErrProjectNotFound):
			pkg.ReturnError(w, http.StatusNotFound, err)
		default:
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		}
		return
	}
	// SSE replay for event-stream clients (progressive rendering), plain
	// JSON otherwise. Framing mirrors StreamLogs.
	if strings.Contains(r.Header.Get("Accept"), "text/event-stream") {
		flusher, ok := w.(http.Flusher)
		if !ok {
			d.logger.Error().Msg("streaming unsupported")
			pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
			return
		}
		writeBuildLogSSE(w, flusher, r, lines)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", lines)
}

// IsBuildRunning implements [DeploymentHandler].
// IsBuildRunning godoc
// @Summary     check the deployment is building new deployment or not.
// @Tags        deployments
// @Produce     json
// @Param       project_name query string true "project name"
// @Param       deployment_name query string true "deployment name"
// @Security    BearerAuth
// @Success     200 {object} pkg.Response{data=schema.IsBuildRunningData} "JSON response"
// @Failure     400 {object} pkg.Response{error=string}
// @Failure     401 {object} pkg.Response{error=string}
// @Failure     500 {object} pkg.Response{error=string}
// @Router      /deployment/build-state [get]
func (d *deploymentHandler) IsBuildRunning(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := pkg.AuthUserIDFromContext(ctx)
	if !ok {
		d.logger.Error().Msg(defined_error.ErrMissingUserIdInContext.Error())
		pkg.ReturnError(w, http.StatusUnauthorized, defined_error.ErrUnauthorized)
		return
	}

	q := r.URL.Query()
	projectName := q.Get("project_name")
	deploymentName := q.Get("deployment_name")
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
		d.logger.Error().Err(defined_error.ErrStringUUIDTypeCasting)
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}

	isBuilding, err := d.ds.IsBuildRunning(ctx, ui, projectName, deploymentName)
	if err != nil {
		d.logger.Error().Err(err).Msg("failed to check deployment build state")
		pkg.ReturnError(w, http.StatusInternalServerError, defined_error.ErrInternalServerError)
		return
	}

	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", schema.IsBuildRunningData{IsBuilding: isBuilding})
}

// writeBuildLogSSE replays build lines as Server-Sent Events: one JSON data
// event per line plus a terminal done event. The stream is bounded by
// sse.buildlog.timeout_seconds so a stuck connection cannot live forever.
func writeBuildLogSSE(w http.ResponseWriter, flusher http.Flusher, r *http.Request, lines []schema.BuildLogLine) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	// Prevent the nginx gateway from buffering the stream.
	w.Header().Set("X-Accel-Buffering", "no")

	ctx, cancel := context.WithTimeout(r.Context(), resolveBuildLogSSETimeout())
	defer cancel()

	for _, line := range lines {
		select {
		case <-ctx.Done():
			return
		default:
		}
		b, _ := json.Marshal(line)
		if _, err := fmt.Fprintf(w, "data: %s\n\n", b); err != nil {
			return
		}
		flusher.Flush()
	}
	if _, err := fmt.Fprintf(w, "event: done\ndata: {\"lines\":%d}\n\n", len(lines)); err != nil {
		return
	}
	flusher.Flush()
}

// resolveBuildLogSSETimeout reads sse.buildlog.timeout_seconds with a
// 300s fallback for missing or non-positive values.
func resolveBuildLogSSETimeout() time.Duration {
	viper.SetDefault("sse.buildlog.timeout_seconds", 300)
	if n := viper.GetInt("sse.buildlog.timeout_seconds"); n > 0 {
		return time.Duration(n) * time.Second
	}
	return 300 * time.Second
}
