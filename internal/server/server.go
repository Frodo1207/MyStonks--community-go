package server

import (
	"MyStonks-go/internal/config"
	user_handler "MyStonks-go/internal/modules/users/handler"
	user_repository "MyStonks-go/internal/modules/users/repository"
	user_service "MyStonks-go/internal/modules/users/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Server struct {
	cfg    *config.Config
	router *gin.Engine
	logger *logrus.Logger
}

func NewServer(cfg *config.Config) *Server {
	// Initialize logger
	logger := logrus.New()
	if cfg.Logger.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}

	// Set gin mode
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create server instance
	s := &Server{
		cfg:    cfg,
		router: gin.New(),
		logger: logger,
	}

	// Configure middleware
	s.router.Use(gin.Recovery())
	s.router.Use(s.loggingMiddleware()) // Now this will work

	// Initialize routes
	s.initializeRoutes()

	return s
}

// loggingMiddleware is now properly defined as a method of Server
func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Before request
		s.logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"ip":     c.ClientIP(),
		}).Info("Incoming request")

		// Process request
		c.Next()

		// After request
		s.logger.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"latency":  c.GetDuration("latency"), // You'd need to set this elsewhere
			"clientIP": c.ClientIP(),
		}).Info("Request completed")
	}
}

func (s *Server) initializeRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Initialize repositories
	userRepo := user_repository.NewUserRepository(s.logger)

	// Initialize services
	userService := user_service.NewUserService(s.logger, userRepo)

	// Initialize handlers
	userHandler := user_handler.NewUserHandler(s.logger, userService)

	// User routes
	userGroup := s.router.Group("/users")
	{
		userGroup.GET("/", userHandler.GetUsers)
		userGroup.GET("/:id", userHandler.GetUser)
		userGroup.POST("/", userHandler.CreateUser)
	}

}

func (s *Server) Run() error {
	s.logger.Infof("Starting %s server on port %d", s.cfg.App.Name, s.cfg.App.Port)
	return s.router.Run(":" + strconv.Itoa(s.cfg.App.Port))
}
