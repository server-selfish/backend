package schema

type (
	GetProjectParams struct {
		ID string `json:"id" validate:"required,min=1"`
	}
	CreateProjectParams struct {
		Name        string `json:"name" validate:"required,min=1"`
		Description string `json:"description"`
	}
	UpdateProjectParams struct {
		ID          string `json:"id" validate:"required,min=1"`
		Name        string `json:"name" validate:"required,min=1"`
		Description string `json:"description"`
	}
	DeleteProjectParams struct {
		ID string `json:"id" validate:"required,min=1"`
	}
)

type (
	GetProjectsData struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)
