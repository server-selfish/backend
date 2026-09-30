package monitoring_infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type (
	// BuildLogLine is a single preserved Docker build output line.
	BuildLogLine struct {
		Time    time.Time
		Seq     int
		Message string
	}
	VictoriaLogsInfra interface {
		QueryBuildLog(ctx context.Context, filter LogFilter) ([]BuildLogLine, error)
	}
	victoriaLogsInfra struct {
		baseURL string
		client  *http.Client
	}
	// LogFilter scopes a build log query. ProjectName + DeploymentName are
	// the primary key: they exist before the deployment row does (a failing
	// first build rolls its transaction back). DeploymentID and AttemptID
	// narrow further when non-empty.
	LogFilter struct {
		DeploymentID   string
		ProjectName    string
		DeploymentName string
		AttemptID      string
		Start          time.Time
		End            time.Time
		Limit          int
	}
)

func NewVictoriaLogsInfra() VictoriaLogsInfra {
	viper.SetDefault("vlogs.query.url", "http://localhost:9428")
	return &victoriaLogsInfra{
		baseURL: viper.GetString("vlogs.query.url"),
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// QueryBuildLog fetches preserved docker build lines from VictoriaLogs,
// oldest first.
func (v *victoriaLogsInfra) QueryBuildLog(ctx context.Context, filter LogFilter) ([]BuildLogLine, error) {
	parts := []string{`log_type:"docker_build"`}
	if filter.DeploymentID != "" {
		parts = append(parts, fmt.Sprintf(`deployment_id:"%s"`, filter.DeploymentID))
	}
	if filter.ProjectName != "" {
		parts = append(parts, fmt.Sprintf(`project_name:"%s"`, filter.ProjectName))
	}
	if filter.DeploymentName != "" {
		parts = append(parts, fmt.Sprintf(`deployment_name:"%s"`, filter.DeploymentName))
	}
	if filter.AttemptID != "" {
		parts = append(parts, fmt.Sprintf(`attempt_id:"%s"`, filter.AttemptID))
	}

	q := url.Values{}
	q.Set("query", strings.Join(parts, " AND "))
	q.Set("start", filter.Start.UTC().Format(time.RFC3339))
	q.Set("end", filter.End.UTC().Format(time.RFC3339))
	q.Set("limit", strconv.Itoa(filter.Limit))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.baseURL+"/select/logsql/query?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("victoria logs query failed with status %d: %.500s", resp.StatusCode, string(body))
	}

	var rows []map[string]any
	// VictoriaLogs streams select results as newline-delimited JSON objects,
	// not a bare array. Tolerate all shapes: array, single object, JSON lines.
	if trimmed := bytes.TrimSpace(body); len(trimmed) > 0 && trimmed[0] == '[' {
		if err := json.Unmarshal(trimmed, &rows); err != nil {
			return nil, err
		}
	} else {
		dec := json.NewDecoder(bytes.NewReader(body))
		for {
			var row map[string]any
			if err := dec.Decode(&row); err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
			rows = append(rows, row)
		}
	}
	lines := make([]BuildLogLine, 0, len(rows))
	for _, row := range rows {
		msg, _ := row["_msg"].(string)
		if msg == "" {
			continue
		}
		var ts time.Time
		if raw, _ := row["_time"].(string); raw != "" {
			ts, _ = time.Parse(time.RFC3339Nano, raw)
		}
		lines = append(lines, BuildLogLine{
			Time:    ts,
			Seq:     logSeq(row["seq"]),
			Message: msg,
		})
	}
	sort.Slice(lines, func(i, j int) bool {
		if lines[i].Time.Equal(lines[j].Time) {
			return lines[i].Seq < lines[j].Seq
		}
		return lines[i].Time.Before(lines[j].Time)
	})
	return lines, nil
}

func logSeq(v any) int {
	switch v := v.(type) {
	case float64:
		return int(v)
	case string:
		n, _ := strconv.Atoi(v)
		return n
	case json.Number:
		n, _ := v.Int64()
		return int(n)
	default:
		return 0
	}
}
