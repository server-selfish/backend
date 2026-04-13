package schema

type (
	CreateDeploymentParams struct {
		Name      string `json:"name" validate:"required,min=1"`
		ProjectID string `json:"project_id" validate:"required,min=1"`
	}
	CreateDeploymentHistoryParams struct {
		DeploymentID string `json:"deployment_id" validate:"required,min=1"`
		GitRemoteUrl string `json:"git_remote_url" validate:"required,min=1"`
		Branch       string `json:"branch" validate:"required,min=1"`
		// CommitID              string  `json:"commit_id" validate:"required,min=1"`
		// CommitMsg             string  `json:"commit_message" validate:"required,min=1"`
		// Version               string  `json:"version" validate:"required,min=1"`
		ExternalPort          []int32 `json:"port"`
		DeploymentTechstackID int32   `json:"techstack_id" validate:"required"`
		BuildCommand          string  `json:"build_command"`
		BuildFolder           string  `json:"build_folder"`
		RunCommand            string  `json:"run_command"`
	}
)
type (
	GetDeploymentData struct {
		DeploymentID      string  `json:"deployment_id"`
		DeploymentName    string  `json:"deployment_name"`
		Branch            string  `json:"branch"`
		CommitID          string  `json:"commit_id"`
		CommitMessage     string  `json:"commit_message"`
		DeploymentVersion string  `json:"deployment_version"`
		Port              []int32 `json:"port"`
		TechstackName     string  `json:"techstack_name"`
		TechstackVersion  string  `json:"techstack_version"`
		ContainerID       string  `json:"container_id"`
		CreatedAt         string  `json:"created_at"`
		UpdatedAt         string  `json:"updated_at"`
	}
	GetSingleDeploymentData struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	GetActiveDeploymentHistory struct {
		DeploymentHistoryID int32   `json:"deployment_history_id"`
		GitRemoteURL        string  `json:"git_remote_url"`
		Branch              string  `json:"branch"`
		CommitId            string  `json:"commit_id"`
		CommitMessage       string  `json:"commit_message"`
		DeploymentVersion   string  `json:"deployment_version"`
		Port                []int32 `json:"port"`
		BuildCommand        string  `json:"build_command"`
		TechstackID         int32   `json:"techstack_id"`
		TechstackName       string  `json:"techstack_name"`
		TechstackVersion    string  `json:"techstack_version"`
	}
	GetHistoryDeploymentHistory struct {
		ID                int32   `json:"deployment_history_id"`
		Branch            string  `json:"branch"`
		CommitID          string  `json:"commit_id"`
		CommitMessage     string  `json:"commit_message"`
		DeploymentVersion string  `json:"deployment_version"`
		Port              []int32 `json:"port"`
		CreatedAt         string  `json:"created_at"`
		UpdatedAt         string  `json:"updated_at"`
	}
)

type (
	DockerFileTemplate struct {
		DockerBaseImage    string
		DockerRuntimeImage string
		BuildCommand       string
		BuildFolder        string
		RunCommand         string
	}
)
