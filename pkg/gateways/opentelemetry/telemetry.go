package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/tccav/identity-service/pkg/config"
)

// InitProvider Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func InitProvider(logger *zap.Logger, version string, config config.Telemetry) (func(), error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.DeploymentEnvironment(config.Environment),
			semconv.ServiceVersion(version),
		),
	)
	if err != nil {
		logger.Error("failed to create oltp resource", zap.Error(err))
		return nil, err
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(config.OtelCollector),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		logger.Error("failed to create trace exporter", zap.Error(err))
		return nil, err
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		err = tracerProvider.Shutdown(ctx)
		if err != nil {
			logger.Error("failed to shutdown TracerProvider", zap.Error(err))
		}

		err = traceExporter.Shutdown(ctx)
		if err != nil {
			logger.Error("failed to shutdown TraceExporter", zap.Error(err))
		}
	}, nil
}
