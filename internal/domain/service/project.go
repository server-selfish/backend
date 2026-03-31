package service

import (
	"context"

	"github.com/server-selfish/backend/internal/domain/dto"
	"github.com/server-selfish/backend/internal/domain/factory"
	"github.com/server-selfish/backend/internal/domain/repository"
	"github.com/server-selfish/backend/internal/pkg"
)

type (
	ProjectService interface {
		CreateProject(ctx context.Context, p *dto.Project) bool
	}
	projectService struct {
		prf      *factory.ProjectRepoFactory
		txRunner pkg.TxRunner
	}
)

func NewProjectService(prf *factory.ProjectRepoFactory) ProjectService {
	return &projectService{
		prf: prf,
	}
}

func (ps *projectService) CreateProject(ctx context.Context, p *dto.Project) bool {
	if err := ps.txRunner.RunInTx(ctx, func(q repository.Querier) error {
		repo := ps.prf.Create(q)
		return repo.CreateProject(p)
	}); err != nil {
		return false
	}
	return true
	// repo := ps.prf.Default()
	// return repo.CreateProject(p)
}
