package service

import (
	"fmt"
	"strings"

	"github.com/server-selfish/backend/internal/constant"
)

func getRunCommandByTechstack(name, mainFilePath, baseImage string) string {
	switch strings.ToLower(name) {
	case "node.js":
		prefix := strings.Split(baseImage, "/")[0]
		if mainFilePath != "" {
			if prefix != "gcr.io" {
				return fmt.Sprintf("node %s", mainFilePath)
			}
			return mainFilePath
		}
		if prefix != "gcr.io" {
			return "node dist/main.js"
		}
		return "main.js"

	case "go":
		if mainFilePath != "" {
			return fmt.Sprintf("/app/%s", mainFilePath)
		}
		return "main"
	case "python":
		if mainFilePath != "" {
			return mainFilePath
		}
		return "main.py"
	default:
		return ""
	}
}

func getFileNameByTechstack(name string) string {
	switch strings.ToLower(name) {
	case "node.js":
		return constant.NODE_DOCKERFILE_TEMPLATE
	case "go":
		return constant.GO_DOCKERFILE_TEMPLATE
	case "python":
		return constant.PYTHON_DOCKERFILE_TEMPLATE
	default:
		return constant.NODE_DOCKERFILE_TEMPLATE
	}
}
