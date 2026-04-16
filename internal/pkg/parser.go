package pkg

import (
	"bytes"
	"html/template"

	docker_infra "github.com/server-selfish/backend/internal/infra/docker"
)

func ParseTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("tmpl").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ParseTemplateFromEmbed(name string, data interface{}) (string, error) {
	content, err := docker_infra.GetTemplate(name)
	if err != nil {
		return "", err
	}

	return ParseTemplate(string(content), data)
}
