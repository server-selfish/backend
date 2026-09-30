package service

import (
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// webhookBuildTimeout bounds a single webhook-triggered build. A Docker build
// plus image pull can run well past the router's 5 minute middleware timeout,
// so each attempt gets its own budget.
const webhookBuildTimeout = 30 * time.Minute

// deployLocks serialises redeploys per deployment id and collapses bursts.
//
// A plain skip-if-busy lock is wrong here: a push is only a notification that
// branch HEAD moved, and the build clones the branch at build time. Dropping a
// push that arrived mid-build would strand the deployment on the older commit
// until the next unrelated push. So a busy deployment records that one more run
// is owed, and the in-flight build repeats once after finishing. Three rapid
// pushes therefore cost two builds, and the second observes the newest HEAD.
type deployLocks struct {
	mu sync.Mutex
	// busy: a build is running for this deployment.
	busy map[pgtype.UUID]bool
	// again: a build finished-or-is-finishing and owes exactly one rerun.
	again map[pgtype.UUID]bool
}

func newDeployLocks() *deployLocks {
	return &deployLocks{
		busy:  make(map[pgtype.UUID]bool),
		again: make(map[pgtype.UUID]bool),
	}
}

// begin claims the lock for id. It reports whether the caller should build now;
// when it returns false the caller has queued a rerun and must not wait, since
// waiting here would block the webhook request that just enqueued it.
func (l *deployLocks) begin(id pgtype.UUID) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.busy[id] {
		// A rerun is already owed, so this push is absorbed by it.
		l.again[id] = true
		return false
	}

	l.busy[id] = true
	return true
}

// done releases the lock and reports whether the caller owes one more run.
func (l *deployLocks) done(id pgtype.UUID) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.again[id] {
		l.again[id] = false
		// Stay busy. The caller rebuilds immediately, and a push arriving in
		// that window must queue behind the rerun. Clearing busy here would let
		// begin hand the lock to a second goroutine, which would build the same
		// deployment concurrently and race on the container name.
		return true
	}

	l.busy[id] = false
	return false
}
