package factory

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/server-selfish/backend/internal/domain/repository"
	storage_infra "github.com/server-selfish/backend/internal/infra/storage"
)

type (
	ProjectRepoFactory struct {
		pgx *pgxpool.Pool
	}
)

func NewProjectRepoFactory(pgx *pgxpool.Pool) *ProjectRepoFactory {
	return &ProjectRepoFactory{pgx: pgx}
}

// Create repo for any Querier (Tx or pool)
func (p *ProjectRepoFactory) Create(q repository.Querier) repository.ProjectRepository {
	return repository.NewProjectRepository(q)
}

// Default repo using pool
func (p *ProjectRepoFactory) Default() repository.ProjectRepository {
	return p.Create(storage_infra.NewQuerier(p.pgx))
}
