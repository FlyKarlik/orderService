package tracer

import (
	"context"
	"fmt"

	"github.com/FlyKarlik/orderService/config"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func New(ctx context.Context, config *config.Config) (func(context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(config.Infrastructure.Opentelemetry.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(fmt.Sprintf(
			"%s:%s",
			config.Infrastructure.Opentelemetry.Host,
			config.Infrastructure.Opentelemetry.Port),
		),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}
