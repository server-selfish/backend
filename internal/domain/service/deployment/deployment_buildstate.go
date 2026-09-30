package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	valkey "github.com/valkey-io/valkey-go"
)

const (
	// buildRunningCachePrefix marks a deployment whose build is in flight. The
	// owner id is part of the key so a caller can only ever read their own
	// namespace: authorisation falls out of the lookup instead of needing a
	// database round trip to prove the deployment belongs to them.
	buildRunningCachePrefix = "deployment_building:"

	// buildRunningMarker is the value stored while a build runs. Existence of
	// the key is the signal, so the value only needs to be non-empty.
	buildRunningMarker = "1"

	// buildRunningTTL must outlast a build. If the key expired mid-run the
	// frontend would report a healthy long build as idle. It is a backstop for
	// the process dying mid-build, not the normal exit path, which deletes the
	// key.
	buildRunningTTL = webhookBuildTimeout + 5*time.Minute
)

// buildRunningCacheKey namespaces a build marker by owner and deployment name.
//
// Names are lowercased because the deployment queries match them with ILIKE,
// while a cache key is an exact match. Normalising both sides keeps a client
// that sends different casing from missing a build that is genuinely running.
func buildRunningCacheKey(userID pgtype.UUID, projectName, deploymentName string) string {
	return fmt.Sprintf(
		"%s%s:%s:%s",
		buildRunningCachePrefix,
		userID.String(),
		strings.ToLower(projectName),
		strings.ToLower(deploymentName),
	)
}

// IsBuildRunning implements [DeploymentService].
//
// A missing key is not an error: nothing is building, which is the common case
// and must not read as a failure.
func (d *deploymentService) IsBuildRunning(
	ctx context.Context,
	userID pgtype.UUID,
	projectName, deploymentName string,
) (bool, error) {
	if _, err := d.cache.Get(ctx, buildRunningCacheKey(userID, projectName, deploymentName)); err != nil {
		if errors.Is(err, valkey.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// markBuildRunning flags a build as in flight before it starts, so a client
// polling straight after the webhook returns already sees it.
//
// Cache failures are logged and swallowed: a missing progress flag must never
// fail a deploy that is otherwise healthy.
func (d *deploymentService) markBuildRunning(
	ctx context.Context,
	userID pgtype.UUID,
	projectName, deploymentName string,
) {
	key := buildRunningCacheKey(userID, projectName, deploymentName)
	if err := d.cache.Set(ctx, key, buildRunningMarker, buildRunningTTL); err != nil {
		d.log.Error().Err(err).Str("key", key).Msg("failed to mark build as running")
	}
}

// clearBuildRunning removes the flag once a build settles, whether it succeeded
// or failed. The build log and the application log carry the outcome, so no
// result state is cached here.
func (d *deploymentService) clearBuildRunning(
	ctx context.Context,
	userID pgtype.UUID,
	projectName, deploymentName string,
) {
	key := buildRunningCacheKey(userID, projectName, deploymentName)
	if err := d.cache.Delete(ctx, key); err != nil {
		d.log.Error().Err(err).Str("key", key).Msg("failed to clear build marker")
	}
}
