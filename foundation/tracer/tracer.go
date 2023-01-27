// Package tracer provides support for distributed tracing.
package tracer

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"log"
)

// Init creates a new trace provider instance and registers it as global trace provider.
func Init(serviceName string, endpoint string, log *log.Logger) error {
	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(endpoint),
		),
	)
	if err != nil {
		return fmt.Errorf("exporting to jaeger: %w", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("service"),
			semconv.DeploymentEnvironmentKey.String("production"),
		)))
	otel.SetTracerProvider(tp)
	return nil
}
