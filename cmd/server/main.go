package main

import (
	"log"
	"os"
	"time"

	"github.com/arvtia/rest-api/internal/config"
	"github.com/arvtia/rest-api/internal/handler"
	"github.com/arvtia/rest-api/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env only in local development
	if os.Getenv("RENDER") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not loaded:", err)
		}
	}

	db := config.InitDB()
	r := gin.Default()

	// CORS setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://your-frontend.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	r.POST("/admin/signup", handler.Signup(db))
	r.POST("/admin/login", handler.Login(db))
	r.POST("/signup", handler.UserSignup(db))
	r.POST("/login", handler.UserLogin(db))

	// Protected admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		admin.POST("/products", handler.CreateProduct(db))
		admin.GET("/products", handler.ListProducts(db))
		admin.PUT("/products/:id", handler.UpdateProduct(db))
		admin.DELETE("/products/:id", handler.DeleteProduct(db))
		admin.POST("/products/:id/media", handler.UploadProductImage(db))
		admin.POST("/products/form", handler.CreateProductWithMedia(db))
	}
	// public route
	r.GET("/products", handler.ListAllProducts(db))

	// Protected user routes
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/details", handler.GetUserDetails(db))
		user.POST("/details", handler.UpdateUserDetails(db))
		// Add cart, orders, payments here later
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
