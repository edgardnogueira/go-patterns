package proxy

import (
	"fmt"
	"sync"
	"time"
)

// MetricType represents the type of metric being tracked.
type MetricType string

const (
	// MetricRequestCount counts the number of requests
	MetricRequestCount MetricType = "request_count"
	// MetricLatency measures operation duration
	MetricLatency MetricType = "latency_ms"
	// MetricErrors counts error occurrences
	MetricErrors MetricType = "error_count"
	// MetricDataSize measures data size processed
	MetricDataSize MetricType = "data_size_bytes"
)

// MetricsProxy collects performance and usage statistics for image operations.
// It provides insights into how images are being accessed and performance characteristics.
type MetricsProxy struct {
	BaseProxy
	metrics     map[string]map[MetricType]float64
	mu          sync.RWMutex
	startTime   time.Time
	metricsHook func(string, MetricType, float64)
}

// NewMetricsProxy creates a new metrics proxy.
func NewMetricsProxy(realImage Image) *MetricsProxy {
	return &MetricsProxy{
		BaseProxy: BaseProxy{
			realImage: realImage,
		},
		metrics:   make(map[string]map[MetricType]float64),
		startTime: time.Now(),
	}
}

// SetMetricsHook sets a callback function that will be called whenever metrics are updated.
// This can be used to integrate with external monitoring systems.
func (p *MetricsProxy) SetMetricsHook(hook func(string, MetricType, float64)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.metricsHook = hook
}

// Display records metrics and delegates to the real image.
func (p *MetricsProxy) Display() error {
	// Track metrics
	p.trackMetric("Display", MetricRequestCount, 1)
	startTime := time.Now()
	
	// Call real operation
	err := p.realImage.Display()
	
	// Record latency and potential error
	latency := float64(time.Since(startTime).Milliseconds())
	p.trackMetric("Display", MetricLatency, latency)
	
	if err != nil {
		p.trackMetric("Display", MetricErrors, 1)
	}
	
	// Track data size
	p.trackMetric("Display", MetricDataSize, float64(p.realImage.GetSize()))
	
	return err
}

// GetFilename records metrics and delegates to the real image.
func (p *MetricsProxy) GetFilename() string {
	p.trackMetric("GetFilename", MetricRequestCount, 1)
	startTime := time.Now()
	
	result := p.realImage.GetFilename()
	
	latency := float64(time.Since(startTime).Milliseconds())
	p.trackMetric("GetFilename", MetricLatency, latency)
	
	return result
}

// GetWidth records metrics and delegates to the real image.
func (p *MetricsProxy) GetWidth() int {
	p.trackMetric("GetWidth", MetricRequestCount, 1)
	startTime := time.Now()
	
	result := p.realImage.GetWidth()
	
	latency := float64(time.Since(startTime).Milliseconds())
	p.trackMetric("GetWidth", MetricLatency, latency)
	
	return result
}

// GetHeight records metrics and delegates to the real image.
func (p *MetricsProxy) GetHeight() int {
	p.trackMetric("GetHeight", MetricRequestCount, 1)
	startTime := time.Now()
	
	result := p.realImage.GetHeight()
	
	latency := float64(time.Since(startTime).Milliseconds())
	p.trackMetric("GetHeight", MetricLatency, latency)
	
	return result
}

// GetSize records metrics and delegates to the real image.
func (p *MetricsProxy) GetSize() int64 {
	p.trackMetric("GetSize", MetricRequestCount, 1)
	startTime := time.Now()
	
	result := p.realImage.GetSize()
	
	latency := float64(time.Since(startTime).Milliseconds())
	p.trackMetric("GetSize", MetricLatency, latency)
	
	return result
}

// GetMetadata records metrics and delegates to the real image.
func (p *MetricsProxy) GetMetadata() map[string]string {
	p.trackMetric("GetMetadata", MetricRequestCount, 1)
	startTime := time.Now()
	
	result := p.realImage.GetMetadata()
	
	latency := float64(time.Since(startTime).Milliseconds())
	p.trackMetric("GetMetadata", MetricLatency, latency)
	
	return result
}

// trackMetric updates a specific metric for an operation.
func (p *MetricsProxy) trackMetric(operation string, metricType MetricType, value float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Initialize operation metrics if needed
	if _, exists := p.metrics[operation]; !exists {
		p.metrics[operation] = make(map[MetricType]float64)
	}
	
	// Update the metric
	if metricType == MetricLatency {
		// For latency, compute running average
		currentCount := p.metrics[operation][MetricRequestCount]
		currentLatency := p.metrics[operation][MetricLatency]
		
		if currentCount > 0 {
			// Compute the new average
			p.metrics[operation][MetricLatency] = (currentLatency*currentCount + value) / (currentCount + 1)
		} else {
			p.metrics[operation][MetricLatency] = value
		}
	} else {
		// For counters, just add the value
		p.metrics[operation][metricType] += value
	}
	
	// Call the metrics hook if set
	if p.metricsHook != nil {
		p.metricsHook(operation, metricType, p.metrics[operation][metricType])
	}
}

// GetMetrics returns all collected metrics.
func (p *MetricsProxy) GetMetrics() map[string]map[MetricType]float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	// Create a deep copy to avoid concurrent access issues
	result := make(map[string]map[MetricType]float64)
	for op, metrics := range p.metrics {
		result[op] = make(map[MetricType]float64)
		for metricType, value := range metrics {
			result[op][metricType] = value
		}
	}
	
	return result
}

// ResetMetrics clears all collected metrics.
func (p *MetricsProxy) ResetMetrics() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.metrics = make(map[string]map[MetricType]float64)
	p.startTime = time.Now()
}

// PrintMetricsReport prints a formatted summary of all metrics.
func (p *MetricsProxy) PrintMetricsReport() {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	fmt.Println("=== Metrics Report ===")
	fmt.Printf("Uptime: %v\n\n", time.Since(p.startTime).Round(time.Second))
	
	for operation, metrics := range p.metrics {
		fmt.Printf("Operation: %s\n", operation)
		fmt.Printf("  Requests: %.0f\n", metrics[MetricRequestCount])
		
		if val, exists := metrics[MetricLatency]; exists {
			fmt.Printf("  Avg Latency: %.2f ms\n", val)
		}
		
		if val, exists := metrics[MetricErrors]; exists && val > 0 {
			fmt.Printf("  Errors: %.0f\n", val)
			errorRate := (val / metrics[MetricRequestCount]) * 100
			fmt.Printf("  Error Rate: %.2f%%\n", errorRate)
		}
		
		if val, exists := metrics[MetricDataSize]; exists {
			fmt.Printf("  Data Size: %.0f bytes\n", val)
		}
		
		fmt.Println()
	}
	fmt.Println("=====================")
}
