package schema

type (
	ContainerStatusResponse struct {
		ContainerStatus string `json:"container_status"`
	}
	ContainerLogEvent struct {
		Stream  string `json:"stream"`
		Message string `json:"message"`
	}
)
