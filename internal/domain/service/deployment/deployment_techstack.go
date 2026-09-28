package service

import (
	"fmt"
	"strings"

	"github.com/server-selfish/backend/internal/constant"
)

func getRunCommandByTechstack(name, buildFolder, mainFileName, runtimeImage string) string {
	switch strings.ToLower(name) {
	case "node.js":
		// Runtime is flattened: /app/{{.BuildFolder}}/* -> /app/*,
		// so CMD must reference the bare file name, not buildFolder/file.
		file := strings.TrimSpace(mainFileName)
		if file == "" {
			file = "index.js"
		}
		// Tolerate "dist/index.js" input: flatten to base name.
		if i := strings.LastIndex(file, "/"); i >= 0 {
			file = file[i+1:]
		}
		if file == "" {
			file = "index.js"
		}
		// Distroless nodejs images use entrypoint /nodejs/bin/node,
		// so CMD is just the file. Shell images need the node prefix.
		// NOTE: must check runtime image, base is always node:*.
		if strings.HasPrefix(runtimeImage, "gcr.io/") {
			return file
		}
		return fmt.Sprintf("node %s", file)

	case "go":
		if strings.TrimSpace(mainFileName) != "" {
			if strings.TrimSpace(buildFolder) != "" {
				return fmt.Sprintf("/app/%s/%s", strings.Trim(buildFolder, "/"), strings.Trim(mainFileName, "/"))
			}
			return fmt.Sprintf("/app/%s", strings.Trim(mainFileName, "/"))
		}
		return "main"
	case "python":
		file := strings.TrimSpace(mainFileName)
		if file == "" {
			file = "main.py"
		}
		// Distroless python images ship a python entrypoint, so CMD is just
		// the file. Full images (e.g. python:*-slim used as runtime) have no
		// entrypoint, so prefix the interpreter explicitly.
		// NOTE: must check runtime image, base is always python:*.
		if strings.HasPrefix(runtimeImage, "gcr.io/") {
			return file
		}
		return fmt.Sprintf("python3 %s", file)
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
