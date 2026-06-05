package schema

import (
	"github.com/jackc/pgx/v5/pgtype"
	deployment_repository "github.com/server-selfish/backend/internal/domain/repository/deployment"
)

type (
	Port struct {
		External int32  `json:"external" validate:"gt=0"`
		Internal int32  `json:"internal" validate:"gt=0"`
		Protocol string `json:"protocol"`
	}
	ENV struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
)
type (
	CreateDeploymentParams struct {
		Name           string `json:"name" validate:"required,min=1"`
		ProjectID      string `json:"project_id" validate:"required,min=1"`
		InstallationID int64  `json:"installation_id" validate:"required"`
		GitRemoteUrl   string `json:"git_remote_url" validate:"required,min=1"`
	}
	CreateDeploymentHistoryParams struct {
		ProjectName           string `json:"project_name" validate:"required,min=1"`
		DeploymentName        string `json:"deployment_name" validate:"required,min=1"`
		DeploymentDescription string `json:"deployment_description"`
		InstallationID        string `json:"installation_id" validate:"required,min=1"`
		RepositoryID          string `json:"repository_id" validate:"required,min=1"`
		Branch                string `json:"branch_name" validate:"required,min=1"`
		DeploymentTechstackID int32  `json:"techstack_id" validate:"required"`
		Env                   []ENV  `json:"env"`
		Port                  []Port `json:"port"`
		BuildCommand          string `json:"build_command"`
		BuildFolder           string `json:"build_folder"`
		MainFileName          string `json:"main_file_name" validate:"required,min=1"`
	}
	BuildAndRunContainerParams struct {
		DepQuery              *deployment_repository.Queries
		ProjectName           string
		Path                  string
		BuildCommand          string
		BuildFolder           string
		DeploymentId          pgtype.UUID
		UserId                pgtype.UUID
		ContainerName         string
		ImageName             string
		Env                   []ENV
		Port                  []Port
		DeploymentTechstackID int32
		MainFileName          string
	}
)
type (
	GetDeploymentData struct {
		DeploymentID      string `json:"deployment_id"`
		GitRemoteUrl      string `json:"git_remote_url"`
		DeploymentName    string `json:"deployment_name"`
		Branch            string `json:"branch"`
		CommitID          string `json:"commit_id"`
		CommitMessage     string `json:"commit_message"`
		DeploymentVersion string `json:"deployment_version"`
		Port              []Port `json:"port"`
		TechstackName     string `json:"techstack_name"`
		TechstackVersion  string `json:"techstack_version"`
		ContainerID       string `json:"container_id"`
		CreatedAt         string `json:"created_at"`
		UpdatedAt         string `json:"updated_at"`
	}
	GetSingleDeploymentData struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		GitRemoteURL string `json:"git_remote_url"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
	}
	GetActiveDeploymentHistory struct {
		DeploymentHistoryID int32  `json:"deployment_history_id"`
		Branch              string `json:"branch"`
		CommitId            string `json:"commit_id"`
		CommitMessage       string `json:"commit_message"`
		DeploymentVersion   string `json:"deployment_version"`
		Port                []Port `json:"port"`
		BuildCommand        string `json:"build_command"`
		TechstackID         int32  `json:"techstack_id"`
		TechstackName       string `json:"techstack_name"`
		TechstackVersion    string `json:"techstack_version"`
	}
	GetHistoryDeploymentHistory struct {
		ID                int32  `json:"deployment_history_id"`
		Branch            string `json:"branch"`
		CommitID          string `json:"commit_id"`
		CommitMessage     string `json:"commit_message"`
		DeploymentVersion string `json:"deployment_version"`
		Port              []Port `json:"port"`
		CreatedAt         string `json:"created_at"`
		UpdatedAt         string `json:"updated_at"`
	}
	GetTechstackList struct {
		Name []string `json:"name"`
	}
	GetTechstackVersion struct {
		ID      int32  `json:"id"`
		Version string `json:"version"`
	}
)

type (
	DockerFileTemplate struct {
		DockerBaseImage    string
		DockerRuntimeImage string
		BuildFolder        string
		BuildCommand       string
		RunCommand         string
	}
)
