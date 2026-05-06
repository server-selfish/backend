package handler

import "github.com/go-chi/chi/v5"

func RegisterPublicAuthRoutes(r chi.Router, ah AuthHandler) {
	r.Get("/auth/github/login", ah.GithubLogin)
	r.Get("/auth/github/callback", ah.GithubCallback)
	r.Post("/auth/refresh", ah.Refresh)
	r.Post("/auth/logout", ah.Logout)
}

func RegisterProtectedAuthRoutes(r chi.Router, ah AuthHandler) {
	r.Get("/auth/me", ah.Me)
}

func RegisterPublicGithubAppRoutes(r chi.Router, gah GithubAppHandler) {
	r.Get("/github-app/callback", gah.Callback)
}

func RegisterProtectedGithubAppRoutes(r chi.Router, gah GithubAppHandler) {
	r.Get("/github-app/install", gah.Install)
	r.Get("/github-app/installations", gah.ListInstallations)
	r.Get("/github-app/installations/{id}/repos", gah.ListInstallationRepositories)
}

func RegisterProjectRoutes(r chi.Router, ph ProjectHandler) {
	r.Post("/project", ph.CreateProject)
	r.Get("/project", ph.GetAllProjects)
	// r.Get("/project/{id}", ph.GetProjectById)
	r.Get("/project/{name}", ph.GetProjectByName)
	r.Patch("/project/{id}", ph.UpdateProjectById)
	r.Delete("/project/{id}", ph.DeleteProjectById)
}

func RegisterDeploymentRoutes(r chi.Router, dh DeploymentHandler) {
	r.Get("/deployment", dh.GetDeploymentsByProjectId)
	r.Get("/deployment/{id}", dh.GetDeploymentByDeploymentId)
	r.Get("/deployment/active/{id}", dh.GetActiveDeploymenByDeploymentId)
	r.Get("/deployment/history/{id}", dh.GetHistoryDeploymentByDeploymentId)
	r.Post("/deployment", dh.CreateDeployment)
	r.Post("/deployment/version", dh.CreateDeployment)
	r.Delete("/deployment/{id}", dh.DeleteDeploymentByDeploymentId)
}
