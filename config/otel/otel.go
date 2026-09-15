// Package otel bootstraps the OpenTelemetry logs pipeline.
//
// Log records are exported over OTLP/HTTP (VictoriaLogs). The constructor
// never fails: when export is disabled via config or the exporter cannot be
// created (e.g. backend unreachable at boot), it returns a stdout-only *Logs
// and records a warning surfaced through Warning(). Runtime export failures
// are absorbed by the batch processor (background retry, drop on overflow)
// and never block the request path.
package otel

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	defaultServiceName    = "selfish-backend"
	defaultServiceVersion = "dev"
	defaultEnvironment    = "local"
	// VictoriaLogs serves OTLP ingestion at this exact path; a bare
	// host:port would fall back to /v1/logs which vlogs does not serve.
	defaultLogsEndpoint = "http://localhost:9428/opentelemetry/api/v1/logs"

	// ScopeName identifies this instrumentation in exported records.
	ScopeName = "github.com/server-selfish/backend/config/logger"

	exportTimeout   = 10 * time.Second
	shutdownTimeout = 5 * time.Second
)

// FlushTimeout bounds the synchronous export used on fatal/panic paths,
// where the process may exit before the batch processor runs.
const FlushTimeout = 5 * time.Second

// Logs owns the OTel logs pipeline. A nil provider means OTLP export is
// disabled or unavailable and the application runs stdout-only.
// All methods are nil-safe.
type Logs struct {
	provider *sdklog.LoggerProvider
	warning  string
}

// NewLogs builds the OTLP exporter, batch processor and logger provider.
// It never returns an error; degradation details are available via Warning().
func NewLogs() *Logs {
	viper.SetDefault("otel.logs.enabled", true)
	viper.SetDefault("otel.logs.endpoint", defaultLogsEndpoint)
	viper.SetDefault("otel.logs.insecure", true)
	viper.SetDefault("otel.service.name", defaultServiceName)
	viper.SetDefault("otel.service.version", defaultServiceVersion)

	if !viper.GetBool("otel.logs.enabled") {
		return &Logs{warning: "otel logs export disabled via otel.logs.enabled, continuing with stdout only"}
	}
	endpoint := viper.GetString("otel.logs.endpoint")

	ctx, cancel := context.WithTimeout(context.Background(), exportTimeout)
	defer cancel()

	opts := []otlploghttp.Option{otlploghttp.WithEndpointURL(endpoint)}
	if viper.GetBool("otel.logs.insecure") {
		opts = append(opts, otlploghttp.WithInsecure())
	}
	exporter, err := otlploghttp.New(ctx, opts...)
	if err != nil {
		return &Logs{warning: "failed to create otlp log exporter for " + endpoint + ": " + err.Error() + ", continuing with stdout only"}
	}

	res, err := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			attribute.String("service.name", firstNonEmpty(viper.GetString("otel.service.name"), defaultServiceName)),
			attribute.String("service.version", firstNonEmpty(viper.GetString("otel.service.version"), defaultServiceVersion)),
			attribute.String("deployment.environment", firstNonEmpty(viper.GetString("ENV"), defaultEnvironment)),
		),
	)
	if err != nil {
		return &Logs{warning: "failed to build otel resource: " + err.Error() + ", continuing with stdout only"}
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)
	return &Logs{provider: provider}
}

// Enabled reports whether OTLP export is active.
func (l *Logs) Enabled() bool {
	return l != nil && l.provider != nil
}

// Warning describes why export is inactive, or "" when healthy.
func (l *Logs) Warning() string {
	if l == nil {
		return ""
	}
	return l.warning
}

// Logger returns an OTel logger bound to the provider's instrumentation
// scope. It must only be called when Enabled() is true.
func (l *Logs) Logger() otellog.Logger {
	return l.provider.Logger(ScopeName)
}

// ForceFlush exports queued records synchronously. Used on fatal/panic
// paths where the process may exit before the batcher runs.
func (l *Logs) ForceFlush(ctx context.Context) error {
	if !l.Enabled() {
		return nil
	}
	return l.provider.ForceFlush(ctx)
}

// Shutdown flushes queued records and releases the exporter.
func (l *Logs) Shutdown() error {
	if !l.Enabled() {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return l.provider.Shutdown(ctx)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
