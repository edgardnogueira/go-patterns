package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	server *http.Server
}

// NewServer creates a new HTTP server
func NewServer() *Server {
	router := gin.Default()
	
	// Add middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Basic health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	return &Server{
		router: router,
	}
}

// Router returns the gin router instance
func (s *Server) Router() *gin.Engine {
	return s.router
}

// Start starts the HTTP server on the specified address
func (s *Server) Start(address string) error {
	s.server = &http.Server{
		Addr:    address,
		Handler: s.router,
	}

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
