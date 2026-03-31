package handler

import (
	"net/http"

	"github.com/server-selfish/backend/internal/domain/service"
)

type (
	ProjectHandler interface {
		CreateProject(w http.ResponseWriter, r *http.Request)
	}
	projectHandler struct {
		ps service.ProjectService
	}
)

func NewProjectHandler(ps service.ProjectService) ProjectHandler {
	return &projectHandler{
		ps: ps,
	}
}

func (p *projectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
