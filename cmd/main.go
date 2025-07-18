package main

import (
	"log"
	"os"

	"movie-api-backend/internal/db"
	"movie-api-backend/internal/handlers"
	"movie-api-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db.InitDB()

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Authentication route
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.GET("/profile", middleware.AuthMiddleware(), handlers.GetProfile)
		}

		// Movie routes
		movies := api.Group("/movies")
		{
			movies.GET("", handlers.GetMovies)
			movies.GET("/:id", handlers.GetMovie)
			movies.GET("/:id/stream", middleware.AuthMiddleware(), handlers.StreamMovie)
			
			// Admin routes
			movies.POST("", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UploadMovie)
			movies.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteMovie)
		}

		// User routes
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/history", handlers.GetViewHistory)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
