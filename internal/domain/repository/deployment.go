package repository

type (
	DeploymentRepository interface{}
	deploymentRepository struct{}
)

func NewDeploymentRepository() DeploymentRepository {
	return &deploymentRepository{}
}
