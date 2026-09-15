package logger

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/server-selfish/backend/config/otel"
	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
)

// otelWriter converts zerolog JSON events into OTel log records.
// Write never fails: unparsable input is dropped so a broken event can
// never break the logging pipeline or the request path.
type otelWriter struct {
	emit  func(otellog.Record)
	flush func(context.Context) error
}

func newOTELWriter(logs *otel.Logs) otelWriter {
	return otelWriter{
		emit: func(r otellog.Record) {
			logs.Logger().Emit(context.Background(), r)
		},
		flush: func(ctx context.Context) error {
			return logs.ForceFlush(ctx)
		},
	}
}

func (w otelWriter) Write(p []byte) (int, error) {
	var evt map[string]any
	if err := json.Unmarshal(p, &evt); err != nil {
		return len(p), nil
	}

	now := time.Now().UTC()
	var rec otellog.Record
	rec.SetObservedTimestamp(now)
	if ts, ok := evt["time"].(string); ok {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			rec.SetTimestamp(t)
		} else {
			rec.SetTimestamp(now)
		}
	} else {
		rec.SetTimestamp(now)
	}

	level, _ := evt["level"].(string)
	rec.SetSeverity(otelSeverity(level))
	rec.SetSeverityText(strings.ToUpper(level))

	if msg, ok := evt["message"].(string); ok {
		rec.SetBody(attribute.StringValue(msg))
	}

	attrs := make([]attribute.KeyValue, 0, len(evt))
	for k, v := range evt {
		if k == "level" || k == "time" || k == "message" {
			continue
		}
		if kv, ok := otelAttr(k, v); ok {
			attrs = append(attrs, kv)
		}
	}
	rec.AddAttributes(attrs...)

	w.emit(rec)

	// Fatal/panic events terminate the process right after Write returns,
	// before the batcher would run: flush synchronously to avoid losing them.
	switch strings.ToLower(level) {
	case "fatal", "panic":
		ctx, cancel := context.WithTimeout(context.Background(), otel.FlushTimeout)
		defer cancel()
		_ = w.flush(ctx)
	}
	return len(p), nil
}

// otelSeverity maps zerolog levels onto OTel severities. OTel has no panic
// severity, so panic maps to fatal with the original text preserved.
func otelSeverity(level string) otellog.Severity {
	switch strings.ToLower(level) {
	case "trace":
		return otellog.SeverityTrace
	case "debug":
		return otellog.SeverityDebug
	case "info":
		return otellog.SeverityInfo
	case "warn", "warning":
		return otellog.SeverityWarn
	case "error":
		return otellog.SeverityError
	case "fatal":
		return otellog.SeverityFatal
	case "panic":
		return otellog.SeverityFatal
	default:
		return otellog.SeverityInfo
	}
}

func otelAttr(k string, v any) (attribute.KeyValue, bool) {
	switch v := v.(type) {
	case string:
		return attribute.String(k, v), true
	case bool:
		return attribute.Bool(k, v), true
	case float64:
		if v == math.Trunc(v) && !math.IsInf(v, 0) && math.Abs(v) < 9e15 {
			return attribute.Int64(k, int64(v)), true
		}
		return attribute.Float64(k, v), true
	case nil:
		return attribute.KeyValue{}, false
	default:
		raw, err := json.Marshal(v)
		if err != nil {
			return attribute.KeyValue{}, false
		}
		return attribute.String(k, string(raw)), true
	}
}
