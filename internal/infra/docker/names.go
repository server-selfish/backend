package docker_infra

import (
	"regexp"
	"strings"
)

var invalidChars = regexp.MustCompile(`[^a-z0-9]+`)
var multipleDashes = regexp.MustCompile(`-+`)

func NormalizeDockerName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = invalidChars.ReplaceAllString(name, "-")
	name = multipleDashes.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")

	if name == "" {
		name = "unnamed"
	}

	return name
}
