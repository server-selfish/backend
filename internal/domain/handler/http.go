package handler

import "github.com/go-chi/chi/v5"

func RegisterProjectRoutes(r *chi.Mux, ph ProjectHandler) {
	r.Post("/project", ph.CreateProject)
	r.Get("/project", ph.GetAllProjects)
	r.Get("/project/{id}", ph.GetProjectById)
	r.Patch("/project/{id}", ph.UpdateProjectById)
	r.Delete("/project/{id}", ph.DeleteProjectById)
}

func RegisterDeploymentRoutes(r *chi.Mux, dh DeploymentHandler) {
	r.Get("/deployment", dh.GetDeploymentsByProjectId)
	r.Get("/deployment/{id}", dh.GetDeploymentByDeploymentId)
	r.Get("/deployment/active/{id}", dh.GetActiveDeploymenByDeploymentId)
	r.Get("/deployment/history/{id}", dh.GetHistoryDeploymentByDeploymentId)
	r.Post("/deployment", dh.CreateDeployment)
	r.Post("/deployment/version", dh.CreateDeployment)
	r.Delete("/deployment/{id}", dh.DeleteDeploymentByDeploymentId)
}
