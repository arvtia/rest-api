package handler

import (
	"net/http"

	"github.com/arvtia/rest-api/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.GetUint("adminID")

		var input model.Product
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		input.AdminID = adminID
		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create product"})
			return
		}

		c.JSON(http.StatusCreated, input)
	}
}

// list products
func ListProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.GetUint("adminID")

		var products []model.Product
		if err := db.Where("admin_id = ?", adminID).Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}

func UpdateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.GetUint("adminID")
		productID := c.Param("id")

		var product model.Product
		if err := db.Where("id = ? AND admin_id = ?", productID, adminID).First(&product).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		var input model.Product
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		db.Model(&product).Updates(input)
		c.JSON(http.StatusOK, product)
	}
}
