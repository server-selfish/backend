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
	}
)
