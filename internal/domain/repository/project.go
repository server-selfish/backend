package repository

import (
	"github.com/server-selfish/backend/internal/domain/dto"
)

type (
	ProjectRepository interface {
		GetProjects() []dto.Project
		GetProject(projectId string) dto.Project
		CreateProject(p *dto.Project) error
		UpdateProject(p *dto.Project) error
		DeleteProjects(projectId string) error
	}
	projectRepository struct {
		q Querier
	}
)

func NewProjectRepository(q Querier) ProjectRepository {
	return &projectRepository{
		q: q,
	}
}

func (pr *projectRepository) CreateProject(p *dto.Project) error {
	panic("unimplemented")
}

// DeleteProjects implements [ProjectRepository].
func (pr *projectRepository) DeleteProjects(projectId string) error {
	panic("unimplemented")
}

// GetProject implements [ProjectRepository].
func (pr *projectRepository) GetProject(projectId string) dto.Project {
	panic("unimplemented")
}

// GetProjects implements [ProjectRepository].
func (pr *projectRepository) GetProjects() []dto.Project {
	panic("unimplemented")
}

// UpdateProject implements [ProjectRepository].
func (pr *projectRepository) UpdateProject(p *dto.Project) error {
	panic("unimplemented")
}
