package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	project_repository "github.com/server-selfish/backend/internal/domain/repository/project"
	"github.com/server-selfish/backend/internal/domain/schema"
	monitoring_infra "github.com/server-selfish/backend/internal/infra/monitoring"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

// ansiEscape strips terminal color codes so preserved lines read cleanly
// in the API and VictoriaLogs UI.
var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")

// buildLogWriter tees Docker build output into structured log records so the
// full build log lands in VictoriaLogs (via the OTLP pipeline) instead of
// being discarded. Write never fails: logging must never break a build.
type buildLogWriter struct {
	mu     sync.Mutex
	buf    strings.Builder
	logger zerolog.Logger
	seq    int
}

func newBuildLogWriter(logger zerolog.Logger, deploymentID pgtype.UUID, projectName, deploymentName, image string) *buildLogWriter {
	depID := "unknown"
	if deploymentID.Valid {
		depID = uuid.UUID(deploymentID.Bytes).String()
	}
	return &buildLogWriter{
		logger: logger.With().
			Str("log_type", "docker_build").
			Str("deployment_id", depID).
			Str("project_name", projectName).
			Str("deployment_name", deploymentName).
			Str("image", image).
			Str("attempt_id", uuid.NewString()).
			Logger(),
	}
}

func (w *buildLogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.buf.Write(p)
	s := w.buf.String()
	idx := strings.LastIndexByte(s, '\n')
	if idx < 0 {
		return len(p), nil
	}
	rest := s[idx+1:]
	w.buf.Reset()
	w.buf.WriteString(rest)
	for _, line := range strings.Split(s[:idx], "\n") {
		w.emit(line)
	}
	return len(p), nil
}

func (w *buildLogWriter) emit(line string) {
	line = ansiEscape.ReplaceAllString(line, "")
	line = strings.TrimRight(line, "\r")
	if strings.TrimSpace(line) == "" {
		return
	}
	w.seq++
	w.logger.Info().Int("seq", w.seq).Msg(line)
}

// finish flushes any trailing partial line and records a summary event so an
// attempt is queryable (and its outcome visible) even when the raw stream is
// empty.
func (w *buildLogWriter) finish(ok bool, elapsed time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if rest := strings.TrimSpace(w.buf.String()); rest != "" {
		w.emit(rest)
	}
	w.buf.Reset()

	status := "success"
	if !ok {
		status = "failed"
	}
	// Sequence after the last line so the summary sorts last among
	// same-timestamp records instead of first.
	w.logger.Info().
		Int("seq", w.seq+1).
		Str("status", status).
		Int("lines", w.seq).
		Dur("duration_ms", elapsed).
		Msg("docker build finished")
}

// GetBuildLog implements [DeploymentService].
func (d *deploymentService) GetBuildLog(ctx context.Context, userId pgtype.UUID, projectName, deploymentName, attemptID string, limit int) ([]schema.BuildLogLine, error) {
	// Ownership check on the project only: the deployment row may not exist
	// yet (the creating build failed and its transaction rolled back), which
	// is exactly when build logs are needed most.
	if _, err := d.pr.GetProjectByName(ctx, project_repository.GetProjectByNameParams{
		UserID: userId,
		Name:   projectName,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, defined_error.ErrProjectNotFound
		}
		return nil, err
	}

	if limit <= 0 {
		limit = 1000
	}
	if limit > 5000 {
		limit = 5000
	}
	end := time.Now()
	start := end.Add(-7 * 24 * time.Hour)

	lines, err := d.vl.QueryBuildLog(ctx, monitoring_infra.LogFilter{
		ProjectName:    projectName,
		DeploymentName: deploymentName,
		AttemptID:      attemptID,
		Start:          start,
		End:            end,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]schema.BuildLogLine, 0, len(lines))
	for _, l := range lines {
		out = append(out, schema.BuildLogLine{
			Time:    l.Time.UTC().Format(time.RFC3339Nano),
			Seq:     l.Seq,
			Message: l.Message,
		})
	}
	return out, nil
}
