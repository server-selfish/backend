package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"

	deployment_repository "github.com/server-selfish/backend/internal/domain/repository/deployment"
	"github.com/server-selfish/backend/internal/domain/schema"
)

var overlapSuffix = regexp.MustCompile(`-\d+$`)

func nextOverlapVersion(prevVersion string, existingCount int64) string {
	base := overlapSuffix.ReplaceAllString(prevVersion, "")
	if base == "" {
		base = prevVersion
	}
	return fmt.Sprintf("%s-%d", base, existingCount+1)
}

func isOnlyDescriptionChange(prev deployment_repository.GetActiveDeploymentDetailByDeploymentNameRow, params schema.UpdateDeploymentParams) (onlyDesc bool, noChange bool) {
	if prev.Branch != params.Branch {
		return false, false
	}
	if prev.TechstackID != params.DeploymentTechstackID {
		return false, false
	}
	if prev.BuildCommand.String != params.BuildCommand {
		return false, false
	}
	if prev.BuildFolder.String != params.BuildFolder {
		return false, false
	}
	if prev.MainFilePath.String != params.MainFileName {
		return false, false
	}
	var prevEnv []schema.ENV
	if err := json.Unmarshal(prev.Env, &prevEnv); err != nil {
		return false, false
	}
	var prevPort []schema.Port
	if err := json.Unmarshal(prev.Port, &prevPort); err != nil {
		return false, false
	}
	nextEnv := slices.Clone(params.Env)
	if nextEnv == nil {
		nextEnv = []schema.ENV{}
	}
	if prevEnv == nil {
		prevEnv = []schema.ENV{}
	}
	nextPort := slices.Clone(params.Port)
	if nextPort == nil {
		nextPort = []schema.Port{}
	}
	if prevPort == nil {
		prevPort = []schema.Port{}
	}
	sortEnv := func(e []schema.ENV) {
		slices.SortFunc(e, func(a, b schema.ENV) int {
			if c := strings.Compare(a.Key, b.Key); c != 0 {
				return c
			}
			return strings.Compare(a.Value, b.Value)
		})
	}
	sortPort := func(p []schema.Port) {
		slices.SortFunc(p, func(a, b schema.Port) int {
			if a.External != b.External {
				if a.External < b.External {
					return -1
				}
				return 1
			}
			if a.Internal != b.Internal {
				if a.Internal < b.Internal {
					return -1
				}
				return 1
			}
			return strings.Compare(a.Protocol, b.Protocol)
		})
	}
	sortEnv(prevEnv)
	sortEnv(nextEnv)
	if !slices.Equal(prevEnv, nextEnv) {
		return false, false
	}
	sortPort(prevPort)
	sortPort(nextPort)
	if !slices.Equal(prevPort, nextPort) {
		return false, false
	}
	if prev.DeploymentDescription.String == params.DeploymentDescription {
		return false, true
	}
	return true, false
}
