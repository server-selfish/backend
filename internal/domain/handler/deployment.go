package handler

type (
	DeploymentHandler interface{}
	deploymentHandler struct{}
)

func NewDeploymentHandler() DeploymentHandler {
	return &deploymentHandler{}
}
