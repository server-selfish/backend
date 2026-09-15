package otel

import (
	"context"
	"testing"
)

func TestLogsNilSafe(t *testing.T) {
	var l *Logs
	if l.Enabled() {
		t.Error("nil Logs must not be enabled")
	}
	if got := l.Warning(); got != "" {
		t.Errorf("nil Logs warning = %q, want empty", got)
	}
	if err := l.Shutdown(); err != nil {
		t.Errorf("nil Logs Shutdown = %v, want nil", err)
	}
	if err := l.ForceFlush(context.Background()); err != nil {
		t.Errorf("nil Logs ForceFlush = %v, want nil", err)
	}
}

func TestDisabledLogs(t *testing.T) {
	l := &Logs{}
	if l.Enabled() {
		t.Error("Logs without provider must not be enabled")
	}
	if err := l.Shutdown(); err != nil {
		t.Errorf("Shutdown = %v, want nil", err)
	}
}
