package main

import (
	"github.com/arvtia/rest-api/internal/config"
	"github.com/arvtia/rest-api/internal/handler"
	"github.com/arvtia/rest-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// initialize DB

	db := config.InitDB()

	//create route
	r := gin.Default()

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
