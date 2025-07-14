package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"boilerplate-service/services"
)

type Server struct {
	logger   *zap.Logger
	services *services.Services
	port     int
	server   *http.Server
	router   *gin.Engine
}

func NewServer(logger *zap.Logger, services *services.Services, port int) *Server {
	return &Server{
		logger:   logger,
		services: services,
		port:     port,
	}
}

func (s *Server) Initialize() {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	s.router = gin.New()

	// Add middleware
	s.router.Use(gin.Recovery())
	s.router.Use(s.loggerMiddleware())
	s.router.Use(s.corsMiddleware())

	// Setup routes
	s.setupRoutes()

	// Create HTTP server
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)
	s.router.GET("/ready", s.readinessCheck)

	// API routes
	api := s.router.Group("/api/v1")
	{
		// Public routes
		api.GET("/status", s.getStatus)

		// Protected routes
		protected := api.Group("")
		protected.Use(s.authMiddleware())
		{
			protected.GET("/profile", s.getProfile)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(s.authMiddleware(), s.adminMiddleware())
		{
			admin.GET("/metrics", s.getMetrics)
		}
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server", zap.Int("port", s.port))
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}

// Middleware
func (s *Server) loggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		s.logger.Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
		)
		return ""
	})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sample auth middleware
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		
		// Set user_id in context for demo
		c.Set("user_id", "sample_user_id")
		c.Next()
	}
}

func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sample admin middleware
		userID := c.GetString("user_id")
		if userID != "admin" {
			c.JSON(403, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Handlers
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "boilerplate-service",
	})
}

func (s *Server) readinessCheck(c *gin.Context) {
	// Check database connection
	if err := s.services.Database.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"error":  "database connection failed",
		})
		return
	}

	// Check Redis connection
	if err := s.services.Redis.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"error":  "redis connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"service": "boilerplate-service",
	})
}

func (s *Server) getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "running",
		"service": "boilerplate-service",
		"version": "1.0.0",
	})
}

func (s *Server) getProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	profile, err := s.services.User.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (s *Server) getMetrics(c *gin.Context) {
	metrics, err := s.services.Admin.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
} 