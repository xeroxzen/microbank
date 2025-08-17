package main

import (
	"log"
	"os"
	"time"

	"microbank/client-service/internal/handlers"
	"microbank/client-service/internal/middleware"
	"microbank/client-service/internal/repository"
	"microbank/client-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection
	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, refreshTokenRepo)
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	adminHandler := handlers.NewAdminHandler(userService)

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	r := gin.Default()

	// Add middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "client-service",
			"timestamp": time.Now().Unix(),
		})
	})

	// Public routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			// Validate token requires authentication
			auth.GET("/validate", middleware.AuthMiddleware(), authHandler.ValidateToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Profile routes
			profile := protected.Group("/profile")
			{
				profile.GET("", userHandler.GetProfile)
				profile.PUT("", userHandler.UpdateProfile)
			}

			// Admin routes - require admin role
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/clients", adminHandler.GetAllClients)
				admin.POST("/clients/:id/blacklist", adminHandler.BlacklistClient)
				admin.DELETE("/clients/:id/blacklist", adminHandler.RemoveFromBlacklist)
			}
		}
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Client Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

