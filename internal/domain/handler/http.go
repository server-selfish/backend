package handler

import "github.com/go-chi/chi/v5"

func RegisterProjectRoutes(r *chi.Mux, ph ProjectHandler) {
	// r.Post("/todo", th.AddNewTodo)
}

func RegisterDeploymentRoutes(r *chi.Mux, dh DeploymentHandler) {
	// r.Post("/todo", th.AddNewTodo)
}
