package main

import (
	"log"
	"os"

	"microbank/banking-service/internal/handlers"
	"microbank/banking-service/internal/middleware"
	"microbank/banking-service/internal/repository"
	"microbank/banking-service/internal/services"

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
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize services
	accountService := services.NewAccountService(accountRepo)
	transactionService := services.NewTransactionService(transactionRepo, accountRepo)

	// Initialize handlers
	accountHandler := handlers.NewAccountHandler(accountService, transactionService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

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
			"status":  "healthy",
			"service": "banking-service",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Protected routes - require authentication
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Account routes
			account := protected.Group("/account")
			{
				account.GET("/balance", accountHandler.GetBalance)
				account.GET("/transactions", accountHandler.GetTransactions)
			}

			// Transaction routes
			transactions := protected.Group("/transactions")
			{
				transactions.POST("/deposit", transactionHandler.Deposit)
				transactions.POST("/withdraw", transactionHandler.Withdraw)
				transactions.GET("/:id", transactionHandler.GetTransaction)
			}
		}
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Banking Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
