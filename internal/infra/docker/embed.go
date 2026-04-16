package docker_infra

import "embed"

//go:embed templates/*
var templatesFS embed.FS

func GetTemplate(name string) (string, error) {
	content, err := templatesFS.ReadFile("templates/" + name + ".tmpl")
	if err != nil {
		return "", err
	}
	return string(content), nil
}
