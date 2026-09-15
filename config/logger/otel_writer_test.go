package logger

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
)

func TestOTelSeverity(t *testing.T) {
	cases := map[string]otellog.Severity{
		"trace":   otellog.SeverityTrace,
		"debug":   otellog.SeverityDebug,
		"info":    otellog.SeverityInfo,
		"warn":    otellog.SeverityWarn,
		"warning": otellog.SeverityWarn,
		"error":   otellog.SeverityError,
		"fatal":   otellog.SeverityFatal,
		"panic":   otellog.SeverityFatal,
		"WARN":    otellog.SeverityWarn,
		"bogus":   otellog.SeverityInfo,
		"":        otellog.SeverityInfo,
	}
	for level, want := range cases {
		if got := otelSeverity(level); got != want {
			t.Errorf("otelSeverity(%q) = %v, want %v", level, got, want)
		}
	}
}

func TestParseLevel(t *testing.T) {
	if got := parseLevel("debug"); got != zerolog.DebugLevel {
		t.Errorf("parseLevel(debug) = %v", got)
	}
	if got := parseLevel("bogus"); got != zerolog.InfoLevel {
		t.Errorf("parseLevel(bogus) = %v, want info fallback", got)
	}
	if got := parseLevel(""); got != zerolog.InfoLevel {
		t.Errorf("parseLevel(\"\") = %v, want info fallback", got)
	}
}

func collectAttrs(rec otellog.Record) map[string]attribute.Value {
	got := map[string]attribute.Value{}
	rec.WalkAttributes(func(kv attribute.KeyValue) bool {
		got[string(kv.Key)] = kv.Value
		return true
	})
	return got
}

func TestOTELWriterEmitsRecord(t *testing.T) {
	var emitted []otellog.Record
	var flushed bool
	w := otelWriter{
		emit:  func(r otellog.Record) { emitted = append(emitted, r) },
		flush: func(context.Context) error { flushed = true; return nil },
	}

	line := `{"level":"error","time":"2026-09-15T02:00:00Z","message":"boom","user":"alice","tries":3,"ratio":1.5,"ok":true,"meta":{"a":1}}`
	n, err := w.Write([]byte(line))
	if err != nil || n != len(line) {
		t.Fatalf("Write = (%d, %v), want (%d, nil)", n, err, len(line))
	}
	if len(emitted) != 1 {
		t.Fatalf("emitted %d records, want 1", len(emitted))
	}
	rec := emitted[0]
	if rec.Severity() != otellog.SeverityError {
		t.Errorf("severity = %v, want error", rec.Severity())
	}
	if rec.SeverityText() != "ERROR" {
		t.Errorf("severity text = %q, want ERROR", rec.SeverityText())
	}
	if rec.Body().AsString() != "boom" {
		t.Errorf("body = %q, want boom", rec.Body().AsString())
	}
	if ts := rec.Timestamp(); ts.Year() != 2026 {
		t.Errorf("timestamp = %v, want 2026 event time", ts)
	}

	attrs := collectAttrs(rec)
	for _, reserved := range []string{"level", "time", "message"} {
		if _, ok := attrs[reserved]; ok {
			t.Errorf("reserved field %q must not become an attribute", reserved)
		}
	}
	if v, ok := attrs["user"]; !ok || v.AsString() != "alice" {
		t.Errorf("user attr = %v", attrs["user"])
	}
	if v, ok := attrs["tries"]; !ok || v.AsInt64() != 3 {
		t.Errorf("tries attr = %v, want int64 3", attrs["tries"])
	}
	if v, ok := attrs["ratio"]; !ok || v.AsFloat64() != 1.5 {
		t.Errorf("ratio attr = %v, want 1.5", attrs["ratio"])
	}
	if v, ok := attrs["ok"]; !ok || !v.AsBool() {
		t.Errorf("ok attr = %v, want true", attrs["ok"])
	}
	if v, ok := attrs["meta"]; !ok || v.AsString() != `{"a":1}` {
		t.Errorf("meta attr = %v, want nested JSON string", attrs["meta"])
	}

	if flushed {
		t.Error("non-fatal event must not flush")
	}
}

func TestOTELWriterDropsMalformedInput(t *testing.T) {
	called := false
	w := otelWriter{
		emit:  func(otellog.Record) { called = true },
		flush: func(context.Context) error { called = true; return nil },
	}
	line := `not json {{{`
	n, err := w.Write([]byte(line))
	if err != nil || n != len(line) {
		t.Fatalf("Write = (%d, %v), want (%d, nil)", n, err, len(line))
	}
	if called {
		t.Error("malformed input must not emit or flush")
	}
}

func TestOTELWriterFlushesOnFatal(t *testing.T) {
	var flushed bool
	w := otelWriter{
		emit:  func(otellog.Record) {},
		flush: func(context.Context) error { flushed = true; return nil },
	}
	if _, err := w.Write([]byte(`{"level":"fatal","message":"dying"}`)); err != nil {
		t.Fatal(err)
	}
	if !flushed {
		t.Error("fatal event must trigger a synchronous flush")
	}
}

func TestOTELWriterTimestampFallback(t *testing.T) {
	var emitted []otellog.Record
	w := otelWriter{
		emit:  func(r otellog.Record) { emitted = append(emitted, r) },
		flush: func(context.Context) error { return nil },
	}
	before := time.Now().UTC()
	if _, err := w.Write([]byte(`{"level":"info","message":"no time field"}`)); err != nil {
		t.Fatal(err)
	}
	if len(emitted) != 1 {
		t.Fatalf("emitted %d records, want 1", len(emitted))
	}
	if ts := emitted[0].Timestamp(); ts.Before(before) || ts.After(time.Now().UTC()) {
		t.Errorf("timestamp = %v, want ~now fallback", ts)
	}
}
