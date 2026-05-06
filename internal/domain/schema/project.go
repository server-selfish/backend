package schema

type (
	CreateProjectParams struct {
		Name        string `json:"name" validate:"required,min=1"`
		Description string `json:"description"`
	}
	UpdateProjectParams struct {
		Name        string `json:"name" validate:"required,min=1"`
		Description string `json:"description"`
	}
)

type (
	GetProjectsData struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
	}
	GetProjectDetail struct {
		ProjectName        string                     `json:"project_name"`
		ProjectDescription string                     `json:"project_description"`
		ProjectCreatedAt   string                     `json:"project_created_at"`
		ProjectUpdatedAt   string                     `json:"project_updated_at"`
		Deployments        []ProjectDeploymentSummary `json:"deployments"`
	}
)

type (
	ProjectDeploymentSummary struct {
		DeploymentName string `json:"deployment_name"`
		TechstackName  string `json:"techstack_name"`
		ContainerName  string `json:"container_name"`
	}
)
