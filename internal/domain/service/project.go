package service

import (
	"context"
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
		CreateProject(ctx context.Context, params *schema.CreateProjectParams) error
		GetAllProjects(ctx context.Context) ([]schema.GetProjectsData, error)
		GetProjectById(ctx context.Context, id pgtype.UUID) (schema.GetProjectsData, error)
		UpdateProjectById(ctx context.Context, id pgtype.UUID, params *schema.UpdateProjectParams) error
		DeleteProjectById(ctx context.Context, id pgtype.UUID) error
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

// DeleteProjectById implements [ProjectService].
func (ps *projectService) DeleteProjectById(ctx context.Context, id pgtype.UUID) error {
	return ps.pr.DeleteProjectById(ctx, id)
}

// GetAllProjects implements [ProjectService].
func (ps *projectService) GetAllProjects(ctx context.Context) ([]schema.GetProjectsData, error) {
	projects, err := ps.pr.GetAllProjects(ctx)
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
		})
	}
	return res, nil
}

// GetProjectById implements [ProjectService].
func (ps *projectService) GetProjectById(ctx context.Context, id pgtype.UUID) (schema.GetProjectsData, error) {
	p, err := ps.pr.GetProjectById(ctx, id)
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
func (ps *projectService) UpdateProjectById(ctx context.Context, id pgtype.UUID, params *schema.UpdateProjectParams) error {
	if err := ps.pr.UpdateProjectById(ctx, project_repository.UpdateProjectByIdParams{
		Name:        params.Name,
		Description: pgtype.Text{String: params.Description},
		ID:          id,
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
func (ps *projectService) CreateProject(ctx context.Context, params *schema.CreateProjectParams) error {
	if err := ps.pr.CreateProject(ctx, project_repository.CreateProjectParams{
		Name:        params.Name,
		Description: pgtype.Text{String: params.Description},
	}); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return pkg.ErrAlreadyExist
		}
		return err
	}
	return nil
}
