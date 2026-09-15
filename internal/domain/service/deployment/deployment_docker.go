package service

import (
	"encoding/json"
	"strings"

	docker_infra "github.com/server-selfish/backend/internal/infra/docker"
	"github.com/server-selfish/backend/internal/pkg"
)

func shellToExecForm(cmd string) string {
	args := strings.Fields(cmd)
	b, _ := json.Marshal(args)
	return string(b)
}

func parseTemplateFromEmbed(name string, data interface{}) (string, error) {
	content, err := docker_infra.GetTemplate(name)
	if err != nil {
		return "", err
	}

	return pkg.ParseTemplate(string(content), data)
}
