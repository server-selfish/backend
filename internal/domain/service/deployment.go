package service

type (
	DeploymentService interface{}
	deploymentService struct{}
)

func NewDeploymentService() DeploymentService {
	return &deploymentService{}
}
