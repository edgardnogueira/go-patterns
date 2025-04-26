package routers

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/common"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/observability"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// SetupRouter configures and returns the router for the API Gateway
func SetupRouter(logger *observability.Logger, orderServiceURL, inventoryServiceURL string) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(loggingMiddleware(logger))
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.RequestID)
	r.Use(observability.TraceMiddleware("api-gateway"))
	r.Use(chimiddleware.Heartbeat("/health"))

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"api-gateway"}`))
	})

	// Order service routes
	r.Route("/api/orders", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Mount("/", proxyHandler(orderServiceURL, "/api/orders"))
	})

	// Inventory service routes
	r.Route("/api/inventory", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Mount("/", proxyHandler(inventoryServiceURL, "/api/inventory"))
	})

	// Public health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API documentation
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/index.html")
	})
	r.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/openapi.json")
	})

	return r
}

// proxyHandler creates a reverse proxy to a backend service
func proxyHandler(targetURL, pathPrefix string) http.Handler {
	target, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add correlation ID for tracing
		if r.Header.Get("X-Request-ID") == "" {
			r.Header.Set("X-Request-ID", common.GenerateRequestID())
		}

		// Forward the request to the target service
		proxy.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs incoming requests
func loggingMiddleware(logger *observability.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			latency := time.Since(start)
			logger.Info("Request processed", map[string]interface{}{
				"method":    r.Method,
				"path":      r.URL.Path,
				"status":    ww.Status(),
				"latency":   latency.String(),
				"size":      ww.BytesWritten(),
				"remote":    r.RemoteAddr,
				"request_id": chimiddleware.GetReqID(r.Context()),
			})
		})
	}
}

// authMiddleware handles authentication and authorization
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In a real-world scenario, you would validate tokens, check permissions, etc.
		// This is a simplified example
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// For demo purposes, we'll just check for a simple token
		if token != "Bearer demo-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), observability.UserIDKey, "demo-user")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
