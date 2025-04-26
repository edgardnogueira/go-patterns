package observability

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// TracingProvider manages the OpenTelemetry tracing configuration
type TracingProvider struct {
	tp *sdktrace.TracerProvider
}

// NewTracingProvider creates a new tracing provider
func NewTracingProvider(serviceName string) (*TracingProvider, error) {
	// Create exporter (using stdout for simplicity, in production use a real exporter)
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("creating stdout exporter: %w", err)
	}

	// Create a resource describing the service
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	// Create the trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set the global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Set the global trace provider
	otel.SetTracerProvider(tp)

	return &TracingProvider{
		tp: tp,
	}, nil
}

// Shutdown stops the trace provider
func (t *TracingProvider) Shutdown(ctx context.Context) error {
	return t.tp.Shutdown(ctx)
}

// Tracer creates a named tracer
func (t *TracingProvider) Tracer(name string) trace.Tracer {
	return t.tp.Tracer(name)
}

// StartSpan starts a new span
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return otel.Tracer("").Start(ctx, name, opts...)
}

// AddAttribute adds an attribute to the current span
func AddAttribute(ctx context.Context, key string, value interface{}) {
	span := trace.SpanFromContext(ctx)
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

// AddEvent adds an event to the current span
func AddEvent(ctx context.Context, name string, attributes ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attributes...))
}

// RecordError records an error in the current span
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
}

// TraceMiddleware creates HTTP middleware for tracing
func TraceMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the context and span from the request
			ctx := r.Context()
			tracer := otel.Tracer(serviceName)
			
			// Start a new span
			ctx, span := tracer.Start(
				ctx,
				fmt.Sprintf("%s %s", r.Method, r.URL.Path),
				trace.WithAttributes(
					semconv.HTTPMethodKey.String(r.Method),
					semconv.HTTPURLKey.String(r.URL.String()),
					semconv.HTTPTargetKey.String(r.URL.Path),
					semconv.HTTPHostKey.String(r.Host),
					semconv.HTTPSchemeKey.String(r.URL.Scheme),
					semconv.HTTPUserAgentKey.String(r.UserAgent()),
				),
			)
			defer span.End()
			
			// Call the next handler with the tracing context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
