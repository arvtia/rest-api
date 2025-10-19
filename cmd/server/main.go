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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://your-frontend.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/admin/signup", handler.Signup(db))
	r.POST("/admin/login", handler.Login(db))

	auth := r.Group("/admin")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/products", handler.CreateProduct(db))
		auth.GET("/products", handler.ListProducts(db))
		auth.PUT("/products/:id", handler.UpdateProduct(db))
		auth.DELETE("/products/:id", handler.DeleteProduct(db))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
