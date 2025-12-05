package tracing

import (
	"adapter/internal/shared/log"
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func InitTracer(ctx context.Context, endpoint, serviceName string) (trace.TracerProvider, error) {
	exporter, err := otlptrace.New(ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		log.Error(context.Background(), err, "Failed to create OTLP exporter")
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("environment", "production"),
		),
	)
	if err != nil {
		log.Error(context.Background(), err, "Failed to create resource")
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

	log.Info(ctx, "OpenTelemetry tracer initialized successfully")
	return tp, nil
}

func StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	tracer := otel.Tracer("adapter-golang-backend")
	return tracer.Start(ctx, spanName, opts...)
}

func StartHTTPSpan(ctx context.Context, r *http.Request, operationName string) (context.Context, trace.Span) {
	tracer := otel.Tracer("adapter-golang-backend")

	opts := []trace.SpanStartOption{
		trace.WithAttributes(
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.URLPath(r.URL.Path),
			semconv.URLScheme(r.URL.Scheme),
			semconv.NetHostName(r.Host),
			semconv.ClientAddress(r.RemoteAddr),
		),
		trace.WithSpanKind(trace.SpanKindServer),
	}

	return tracer.Start(ctx, operationName, opts...)
}

func StartBusinessSpan(ctx context.Context, operation string, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	tracer := otel.Tracer("adapter-golang-backend")

	attrs := []attribute.KeyValue{
		attribute.String("operation.type", "business"),
		attribute.String("operation.name", operation),
	}
	attrs = append(attrs, attributes...)

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindInternal),
	}

	return tracer.Start(ctx, operation, opts...)
}

func RecordError(span trace.Span, err error, description string) {
	if err == nil {
		return
	}

	span.RecordError(err, trace.WithAttributes(
		attribute.String("error.type", fmt.Sprintf("%T", err)),
		attribute.String("error.description", description),
	))
	span.SetStatus(codes.Error, description)
}

func SetHTTPStatus(span trace.Span, statusCode int) {
	span.SetAttributes(semconv.HTTPResponseStatusCode(statusCode))

	if statusCode >= 400 {
		span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", statusCode))
	} else {
		span.SetStatus(codes.Ok, "")
	}
}

func AddBusinessAttributes(span trace.Span, attrs map[string]interface{}) {
	for key, value := range attrs {
		switch v := value.(type) {
		case string:
			span.SetAttributes(attribute.String(key, v))
		case int:
			span.SetAttributes(attribute.Int(key, v))
		case int64:
			span.SetAttributes(attribute.Int64(key, v))
		case float64:
			span.SetAttributes(attribute.Float64(key, v))
		case bool:
			span.SetAttributes(attribute.Bool(key, v))
		default:
			span.SetAttributes(attribute.String(key, fmt.Sprintf("%v", v)))
		}
	}
}
