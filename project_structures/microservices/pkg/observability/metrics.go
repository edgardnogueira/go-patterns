package observability

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsServer provides prometheus metrics functionality
type MetricsServer struct {
	registry  *prometheus.Registry
	server    *http.Server
	Counters  map[string]prometheus.Counter
	Gauges    map[string]prometheus.Gauge
	Histograms map[string]prometheus.Histogram
	Summaries map[string]prometheus.Summary
}

// NewMetricsServer creates a new metrics server
func NewMetricsServer(serviceName string, port int) *MetricsServer {
	registry := prometheus.NewRegistry()
	
	// Register the Go collectors
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	
	ms := &MetricsServer{
		registry:   registry,
		Counters:   make(map[string]prometheus.Counter),
		Gauges:     make(map[string]prometheus.Gauge),
		Histograms: make(map[string]prometheus.Histogram),
		Summaries:  make(map[string]prometheus.Summary),
		server: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
	}

	// Create a handler for metrics
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	
	return ms
}

// Start starts the metrics server
func (ms *MetricsServer) Start() error {
	go func() {
		if err := ms.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Failed to start metrics server: %v", err))
		}
	}()
	return nil
}

// Stop stops the metrics server
func (ms *MetricsServer) Stop() error {
	return ms.server.Close()
}

// NewCounter creates a new counter metric
func (ms *MetricsServer) NewCounter(name, help string, labels ...string) prometheus.Counter {
	counter := promauto.With(ms.registry).NewCounter(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
	)
	ms.Counters[name] = counter
	return counter
}

// NewGauge creates a new gauge metric
func (ms *MetricsServer) NewGauge(name, help string, labels ...string) prometheus.Gauge {
	gauge := promauto.With(ms.registry).NewGauge(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
	)
	ms.Gauges[name] = gauge
	return gauge
}

// NewHistogram creates a new histogram metric
func (ms *MetricsServer) NewHistogram(name, help string, buckets []float64, labels ...string) prometheus.Histogram {
	histogram := promauto.With(ms.registry).NewHistogram(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    help,
			Buckets: buckets,
		},
	)
	ms.Histograms[name] = histogram
	return histogram
}

// NewSummary creates a new summary metric
func (ms *MetricsServer) NewSummary(name, help string, objectives map[float64]float64, labels ...string) prometheus.Summary {
	summary := promauto.With(ms.registry).NewSummary(
		prometheus.SummaryOpts{
			Name:       name,
			Help:       help,
			Objectives: objectives,
		},
	)
	ms.Summaries[name] = summary
	return summary
}

// Timer is a utility for timing code execution
type Timer struct {
	histogram prometheus.Histogram
	startTime time.Time
}

// NewTimer creates a new timer
func NewTimer(histogram prometheus.Histogram) *Timer {
	return &Timer{
		histogram: histogram,
		startTime: time.Now(),
	}
}

// ObserveDuration records the duration
func (t *Timer) ObserveDuration() {
	duration := time.Since(t.startTime).Seconds()
	t.histogram.Observe(duration)
}

// Middleware creates a middleware for recording HTTP request metrics
func (ms *MetricsServer) Middleware(serviceName string) func(http.Handler) http.Handler {
	// Create metrics
	requestCounter := ms.NewCounter(
		fmt.Sprintf("%s_http_requests_total", serviceName),
		"Total number of HTTP requests",
	)
	
	requestDuration := ms.NewHistogram(
		fmt.Sprintf("%s_http_request_duration_seconds", serviceName),
		"HTTP request duration in seconds",
		prometheus.DefBuckets,
	)

	// Return the middleware function
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Start the timer
			timer := NewTimer(requestDuration)
			
			// Call the next handler
			next.ServeHTTP(w, r)
			
			// Increment request counter
			requestCounter.Inc()
			
			// Record request duration
			timer.ObserveDuration()
		})
	}
}
