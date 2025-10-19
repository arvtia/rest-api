package main

import (
	"log"
	"time"

	"github.com/arvtia/rest-api/internal/config"
	"github.com/arvtia/rest-api/internal/handler"
	"github.com/arvtia/rest-api/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// initialize DB

	// Load .env file before anything else
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db := config.InitDB()

	//create route
	r := gin.Default()

	// CORS config
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://your-frontend.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))


	// public routes
	r.POST("/admin/signup", handler.Signup(db))
	r.POST("/admin/login", handler.Login(db))

	// protected Routes
	auth := r.Group("/admin")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/products", handler.CreateProduct(db))
		auth.GET("/products", handler.ListProducts(db))
		auth.PUT("/products/:id", handler.UpdateProduct(db))
		auth.DELETE("products/:id", handler.DeleteProduct(db))
	}

	r.Run(":8080") // or env config
}
