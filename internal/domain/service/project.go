package service

import (
	"context"

	"github.com/server-selfish/backend/internal/domain/dto"
	project_repository "github.com/server-selfish/backend/internal/domain/repository/project"
)

type (
	ProjectService interface {
		CreateProject(ctx context.Context, p *dto.Project) bool
	}
	projectService struct {
		pr *project_repository.Queries
	}
)

func NewProjectService(pr *project_repository.Queries) ProjectService {
	return &projectService{
		pr: pr,
	}
}

func (ps *projectService) CreateProject(ctx context.Context, p *dto.Project) bool {
	// if err := ps.txRunner.RunInTx(ctx, func(q repository.Querier) error {
	// 	repo := ps.prf.Create(q)
	// 	return repo.CreateProject(p)
	// }); err != nil {
	// 	return false
	// }
	// return true
	// repo := ps.prf.Default()
	// return repo.CreateProject(p)
	panic("undimplemented")
}
