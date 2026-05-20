package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	deployment_repository "github.com/server-selfish/backend/internal/domain/repository/deployment"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/domain/service"
	"github.com/server-selfish/backend/internal/pkg"
)

type (
	DeploymentHandler interface {
		GetDeploymentsByProjectId(w http.ResponseWriter, r *http.Request)
		GetDeploymentByDeploymentId(w http.ResponseWriter, r *http.Request)
		GetActiveDeploymenByDeploymentId(w http.ResponseWriter, r *http.Request)
		GetHistoryDeploymentByDeploymentId(w http.ResponseWriter, r *http.Request)
		GetTechstackName(w http.ResponseWriter, r *http.Request)
		GetTechstackVersionByName(w http.ResponseWriter, r *http.Request)
		CreateDeployment(w http.ResponseWriter, r *http.Request)
		CreateNewDeploymentVersionByDeploymentId(w http.ResponseWriter, r *http.Request)
		DeleteDeploymentByDeploymentId(w http.ResponseWriter, r *http.Request)
	}
	deploymentHandler struct {
		ds service.DeploymentService
	}
)

func NewDeploymentHandler(ds service.DeploymentService) DeploymentHandler {
	return &deploymentHandler{
		ds: ds,
	}
}

// GetTechstackName implements [DeploymentHandler].
func (d *deploymentHandler) GetTechstackName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tl, err := d.ds.GetTechstackName(ctx)
	if err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch techstack sucess", tl)
}

// GetTechstackVersionByName implements [DeploymentHandler].
func (d *deploymentHandler) GetTechstackVersionByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tn := chi.URLParam(r, "techstack_name")
	if tn == "" {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	vl, err := d.ds.GetTechstackVersionByName(ctx, tn)
	if err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch version success", vl)
}

// CreateNewDeploymentVersion implements [DeploymentHandler].
func (d *deploymentHandler) CreateNewDeploymentVersionByDeploymentId(w http.ResponseWriter, r *http.Request) {
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
	req, ok := pkg.DecodeAndValidateBody[schema.CreateDeploymentHistoryParams](w, r)
	if !ok {
		return
	}
	id, err := pkg.StringToPgUUID(req.DeploymentID)
	if err != nil {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	ii, err := strconv.ParseInt(req.InstallationID, 10, 64)
	if err != nil {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	if err := d.ds.CreateNewDeploymentVersionByDeploymentId(ctx, ui, ii, deployment_repository.CreateDeploymentHistoryParams{
		DeploymentID:          id,
		Branch:                req.Branch,
		ExternalPort:          req.ExternalPort,
		DeploymentTechstackID: req.DeploymentTechstackID,
		BuildCommand:          pgtype.Text{String: req.BuildCommand},
		BuildFolder:           pgtype.Text{String: req.BuildFolder},
		RunCommand:            pgtype.Text{String: req.RunCommand},
	}); err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "new version deployed", nil)
}

// CreateDeployment implements [DeploymentHandler].
func (d *deploymentHandler) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	req, ok := pkg.DecodeAndValidateBody[schema.CreateDeploymentParams](w, r)
	if !ok {
		return
	}
	ctx := r.Context()
	id, err := pkg.StringToPgUUID(req.ProjectID)
	if err != nil {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	if err := d.ds.CreateDeployment(ctx, deployment_repository.CreateDeploymentParams{
		Name:           req.Name,
		ProjectID:      id,
		GitRemoteUrl:   req.GitRemoteUrl,
		InstallationID: req.InstallationID,
	}); err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "deployment created", nil)
}

// DeleteDeploymentByDeploymentId implements [DeploymentHandler].
func (d *deploymentHandler) DeleteDeploymentByDeploymentId(w http.ResponseWriter, r *http.Request) {
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
	if err := d.ds.DeleteDeploymentByDeploymentId(ctx, ui, id); err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "deployment deleted", nil)
}

// GetHistoryDeploymentByDeploymentId implements [DeploymentHandler].
func (d *deploymentHandler) GetHistoryDeploymentByDeploymentId(w http.ResponseWriter, r *http.Request) {
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
	deployments, err := d.ds.GetHistoryDeploymentByDeploymentId(ctx, ui, id)
	if err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", deployments)
}

// GetActiveDeploymenByDeploymentId implements [DeploymentHandler].
func (d *deploymentHandler) GetActiveDeploymenByDeploymentId(w http.ResponseWriter, r *http.Request) {
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

	deployment, err := d.ds.GetActiveDeploymentByDeploymentId(ctx, ui, id)
	if err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", deployment)
}

// GetDeploymentByDeploymentId implements [DeploymentHandler].
func (d *deploymentHandler) GetDeploymentByDeploymentId(w http.ResponseWriter, r *http.Request) {
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
	deployment, err := d.ds.GetDeploymentByDeploymentId(ctx, ui, id)
	if err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", deployment)
}

// GetDeploymentsByProjectId implements [DeploymentHandler].
func (d *deploymentHandler) GetDeploymentsByProjectId(w http.ResponseWriter, r *http.Request) {
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

	q := r.URL.Query()
	projectId := q.Get("project_id")
	if projectId == "" {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	id, err := pkg.StringToPgUUID(projectId)
	if err != nil {
		pkg.ReturnError(w, pkg.ErrBadRequest)
		return
	}
	deployments, err := d.ds.GetDeploymentsByProjectId(ctx, ui, id)
	if err != nil {
		pkg.ReturnError(w, err)
		return
	}
	pkg.ReturnSuccess(w, http.StatusOK, "fetch success", deployments)
}
