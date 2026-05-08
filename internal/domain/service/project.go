package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	project_repository "github.com/server-selfish/backend/internal/domain/repository/project"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/pkg"
)

type (
	ProjectService interface {
		CreateProject(ctx context.Context, userId pgtype.UUID, params *schema.CreateProjectParams) error
		GetAllProjects(ctx context.Context, userId pgtype.UUID) ([]schema.GetProjectsData, error)
		GetProjectById(ctx context.Context, id, userID pgtype.UUID) (schema.GetProjectsData, error)
		GetProjectByName(ctx context.Context, name string, userID pgtype.UUID) (schema.GetProjectDetail, error)
		GetProjectByNameDetail(ctx context.Context, name string, userID pgtype.UUID) (schema.GetProjectAllDetail, error)
		UpdateProjectById(ctx context.Context, id, userID pgtype.UUID, params *schema.UpdateProjectParams) error
		DeleteProjectById(ctx context.Context, id, userID pgtype.UUID) error
	}
	projectService struct {
		pr *project_repository.Queries
		tm pkg.TxManager
	}
)

func NewProjectService(pr *project_repository.Queries, tm pkg.TxManager) ProjectService {
	return &projectService{
		pr: pr,
		tm: tm,
	}
}

// GetProjectByNameDetail implements [ProjectService].
func (ps *projectService) GetProjectByNameDetail(ctx context.Context, name string, userID pgtype.UUID) (schema.GetProjectAllDetail, error) {
	p, err := ps.pr.GetProjectByNameDetail(ctx, project_repository.GetProjectByNameDetailParams{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return schema.GetProjectAllDetail{}, pkg.ErrNotFound
		}
		return schema.GetProjectAllDetail{}, err
	}
	var deployments []schema.ProjectDeploymentDetail
	if err := json.Unmarshal(p.Deployments, &deployments); err != nil {
		return schema.GetProjectAllDetail{}, err
	}
	res := schema.GetProjectAllDetail{
		ProjectName:        p.ProjectName,
		ProjectDescription: p.ProjectDescription.String,
		ProjectCreatedAt:   p.ProjectCreatedAt.Time.String(),
		ProjectUpdatedAt:   p.ProjectUpdatedAt.Time.String(),
		Deployments:        deployments,
	}
	return res, nil
}

// GetProjectByName implements [ProjectService].
func (ps *projectService) GetProjectByName(ctx context.Context, name string, userID pgtype.UUID) (schema.GetProjectDetail, error) {
	p, err := ps.pr.GetProjectByName(ctx, project_repository.GetProjectByNameParams{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return schema.GetProjectDetail{}, pkg.ErrNotFound
		}
		return schema.GetProjectDetail{}, err
	}
	var deployments []schema.ProjectDeploymentSummary
	if err := json.Unmarshal(p.Deployments, &deployments); err != nil {
		return schema.GetProjectDetail{}, err
	}
	res := schema.GetProjectDetail{
		ProjectName:        p.ProjectName,
		ProjectDescription: p.ProjectDescription.String,
		ProjectCreatedAt:   p.ProjectCreatedAt.Time.String(),
		ProjectUpdatedAt:   p.ProjectUpdatedAt.Time.String(),
		Deployments:        deployments,
	}
	return res, nil
}

// DeleteProjectById implements [ProjectService].
func (ps *projectService) DeleteProjectById(ctx context.Context, id, userID pgtype.UUID) error {
	return ps.pr.DeleteProjectById(ctx, project_repository.DeleteProjectByIdParams{
		ID:     id,
		UserID: userID,
	})
}

// GetAllProjects implements [ProjectService].
func (ps *projectService) GetAllProjects(ctx context.Context, userID pgtype.UUID) ([]schema.GetProjectsData, error) {
	projects, err := ps.pr.GetAllProjects(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(projects) == 0 {
		return nil, pkg.ErrNotFound
	}
	var res []schema.GetProjectsData
	for _, p := range projects {
		res = append(res, schema.GetProjectsData{
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description.String,
			CreatedAt:   p.CreatedAt.Time.String(),
		})
	}
	return res, nil
}

// GetProjectById implements [ProjectService].
func (ps *projectService) GetProjectById(ctx context.Context, id, userID pgtype.UUID) (schema.GetProjectsData, error) {
	p, err := ps.pr.GetProjectById(ctx, project_repository.GetProjectByIdParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return schema.GetProjectsData{}, pkg.ErrNotFound
		}
		return schema.GetProjectsData{}, err
	}
	res := schema.GetProjectsData{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description.String,
	}
	return res, nil
}

// UpdateProjectById implements [ProjectService].
func (ps *projectService) UpdateProjectById(ctx context.Context, id, userID pgtype.UUID, params *schema.UpdateProjectParams) error {
	if err := ps.pr.UpdateProjectById(ctx, project_repository.UpdateProjectByIdParams{
		Name:        params.Name,
		Description: pgtype.Text{String: params.Description},
		ID:          id,
		UserID:      userID,
	}); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return pkg.ErrAlreadyExist
		}
		return err
	}
	return nil
}

// UpdateProjectById implements [ProjectService].
func (ps *projectService) CreateProject(ctx context.Context, userID pgtype.UUID, params *schema.CreateProjectParams) error {
	if err := ps.pr.CreateProject(ctx, project_repository.CreateProjectParams{
		Name:        params.Name,
		Description: pgtype.Text{String: params.Description, Valid: true},
		UserID:      userID,
	}); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return pkg.ErrAlreadyExist
		}
		return err
	}
	return nil
}
