package otel_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/server-selfish/backend/config/logger"
	"github.com/server-selfish/backend/config/otel"
	"github.com/spf13/viper"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	"google.golang.org/protobuf/proto"
)

type capturedRequest struct {
	path        string
	contentType string
	req         *collogspb.ExportLogsServiceRequest
}

// TestOTLPWireExport spins a stub OTLP/HTTP server, runs the real provider +
// zerolog bridge against it, and validates the protobuf payload on the wire.
func TestOTLPWireExport(t *testing.T) {
	var mu sync.Mutex
	var got []capturedRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
		}
		var req collogspb.ExportLogsServiceRequest
		if err := proto.Unmarshal(body, &req); err != nil {
			t.Errorf("unmarshal OTLP payload: %v", err)
		}
		mu.Lock()
		got = append(got, capturedRequest{
			path:        r.URL.Path,
			contentType: r.Header.Get("Content-Type"),
			req:         &req,
		})
		mu.Unlock()
		resp, _ := proto.Marshal(&collogspb.ExportLogsServiceResponse{})
		w.Header().Set("Content-Type", "application/x-protobuf")
		_, _ = w.Write(resp)
	}))
	defer srv.Close()

	const marker = "otlp-wire-smoke-9f3a"
	setViper(t, "otel.logs.enabled", true)
	setViper(t, "otel.logs.endpoint", srv.URL+"/opentelemetry/api/v1/logs")
	setViper(t, "otel.logs.insecure", true)
	setViper(t, "otel.service.name", "smoke-backend")
	setViper(t, "otel.service.version", "test")
	setViper(t, "log.stdout.enabled", false)

	logs := otel.NewLogs()
	if !logs.Enabled() {
		t.Fatalf("provider not enabled: %q", logs.Warning())
	}
	appLogger := logger.NewLogger(logs)
	appLogger.Error().Str("user", "alice").Int("tries", 3).Msg(marker)
	if err := logs.Shutdown(); err != nil {
		t.Fatalf("shutdown: %v", err)
	}

	if len(got) == 0 {
		t.Fatal("stub server received no export requests")
	}
	for _, c := range got {
		if c.path != "/opentelemetry/api/v1/logs" {
			t.Errorf("request path = %q", c.path)
		}
		if c.contentType != "application/x-protobuf" {
			t.Errorf("content-type = %q", c.contentType)
		}
	}

	var found bool
	for _, c := range got {
		for _, rl := range c.req.ResourceLogs {
			if strAttr(rl.Resource.Attributes, "service.name") != "smoke-backend" {
				t.Errorf("service.name = %q, want smoke-backend", strAttr(rl.Resource.Attributes, "service.name"))
			}
			for _, sl := range rl.ScopeLogs {
				if sl.Scope.GetName() != otel.ScopeName {
					t.Errorf("scope = %q, want %q", sl.Scope.GetName(), otel.ScopeName)
				}
				for _, rec := range sl.LogRecords {
					if rec.Body.GetStringValue() != marker {
						continue
					}
					found = true
					if rec.SeverityNumber != logspb.SeverityNumber_SEVERITY_NUMBER_ERROR {
						t.Errorf("severity = %v, want ERROR", rec.SeverityNumber)
					}
					if rec.SeverityText != "ERROR" {
						t.Errorf("severity text = %q", rec.SeverityText)
					}
					if strAttr(rec.Attributes, "user") != "alice" {
						t.Errorf("user attr missing: %v", rec.Attributes)
					}
					if intAttr(rec.Attributes, "tries") != 3 {
						t.Errorf("tries attr missing: %v", rec.Attributes)
					}
				}
			}
		}
	}
	if !found {
		t.Errorf("marker record %q not found in %d export requests", marker, len(got))
	}
}

// TestDisabledExporterBoots ensures the app can always construct a logger,
// even with export explicitly off.
func TestDisabledExporterBoots(t *testing.T) {
	setViper(t, "otel.logs.enabled", false)
	logs := otel.NewLogs()
	if logs.Enabled() {
		t.Error("provider must be disabled")
	}
	if logs.Warning() == "" {
		t.Error("expected a degradation warning")
	}
	appLogger := logger.NewLogger(logs)
	appLogger.Info().Msg("stdout only")
	if err := logs.Shutdown(); err != nil {
		t.Errorf("shutdown = %v", err)
	}
}

func setViper(t *testing.T, key string, value any) {
	t.Helper()
	prev := viper.Get(key)
	viper.Set(key, value)
	t.Cleanup(func() { viper.Set(key, prev) })
}

func strAttr(attrs []*commonv1.KeyValue, key string) string {
	for _, kv := range attrs {
		if kv.Key == key {
			return kv.Value.GetStringValue()
		}
	}
	return ""
}

func intAttr(attrs []*commonv1.KeyValue, key string) int64 {
	for _, kv := range attrs {
		if kv.Key == key {
			return kv.Value.GetIntValue()
		}
	}
	return -1
}
