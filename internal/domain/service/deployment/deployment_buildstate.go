package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
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

	// deployLockCachePrefix guards once-at-a-time builds per deployment.
	// Separate from buildRunningCachePrefix: progress flag uses plain SET
	// and allows overwrite, mutex uses SET NX and rejects second owner.
	deployLockCachePrefix = "deploy_lock:"
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

// deployLockCacheKey namespaces a mutex by owner and deployment name.
// Same lowercasing as buildRunningCacheKey so manual and webhook paths
// contend on one key regardless of casing.
func deployLockCacheKey(userID pgtype.UUID, projectName, deploymentName string) string {
	return fmt.Sprintf(
		"%s%s:%s:%s",
		deployLockCachePrefix,
		userID.String(),
		strings.ToLower(projectName),
		strings.ToLower(deploymentName),
	)
}

// tryAcquireDeployLock claims once-at-a-time ownership via SET NX.
// Returns defined_error.ErrBuildInProgress when another build owns key.
// Cache errors fail closed: duplicate container builds risk worse than
// a rejected request during cache outage.
// ponytail: simple Delete release, no owner token. Add token + Lua
// compare-del when lock steal observed.
func (d *deploymentService) tryAcquireDeployLock(
	ctx context.Context,
	userID pgtype.UUID,
	projectName, deploymentName string,
) error {
	key := deployLockCacheKey(userID, projectName, deploymentName)
	acquired, err := d.cache.TryAcquire(ctx, key, buildRunningMarker, buildRunningTTL)
	if err != nil {
		d.log.Error().Err(err).Str("key", key).Msg("failed to acquire deploy lock")
		return err
	}
	if !acquired {
		return defined_error.ErrBuildInProgress
	}
	return nil
}

// releaseDeployLock frees mutex. Caller passes detached context so
// release survives HTTP or build timeout cancellation.
func (d *deploymentService) releaseDeployLock(
	ctx context.Context,
	userID pgtype.UUID,
	projectName, deploymentName string,
) {
	key := deployLockCacheKey(userID, projectName, deploymentName)
	if err := d.cache.Delete(ctx, key); err != nil {
		d.log.Error().Err(err).Str("key", key).Msg("failed to release deploy lock")
	}
}
