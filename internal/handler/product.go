package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/arvtia/rest-api/internal/model"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// create products
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

// upload with images
func CreateProductWithMedia(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.GetUint("adminID")

		// Parse form fields with validation
		name := c.PostForm("name")
		description := c.PostForm("description")
		priceStr := c.PostForm("price")
		stockStr := c.PostForm("stock")
		category := c.PostForm("category")

		if name == "" || priceStr == "" || stockStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
			return
		}

		stock, err := strconv.Atoi(stockStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock format"})
			return
		}

		product := model.Product{
			Name:        name,
			Description: description,
			Price:       price,
			Stock:       stock,
			Category:    category,
			AdminID:     adminID,
		}

		if err := db.Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create product"})
			return
		}

		// Handle multiple image uploads
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		files := form.File["images"]
		if len(files) > 0 {
			cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary config failed"})
				return
			}

			var mediaRecords []model.ProductMedia
			for i, file := range files {
				src, err := file.Open()
				if err != nil {
					continue
				}
				defer src.Close()

				uploadResp, err := cld.Upload.Upload(context.Background(), src, uploader.UploadParams{
					Folder:   "products",
					PublicID: fmt.Sprintf("product_%d_img_%d", product.ID, i),
				})
				if err != nil {
					continue
				}

				mediaRecords = append(mediaRecords, model.ProductMedia{
					ProductID: product.ID,
					URL:       uploadResp.SecureURL,
					Type:      "image",
					AltText:   file.Filename,
				})
			}

			if len(mediaRecords) > 0 {
				if err := db.Create(&mediaRecords).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save media"})
					return
				}
				// Attach media to product for response
				product.Media = mediaRecords
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"product": product,
		})
	}
}

// fething the products list
func ListProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.GetUint("adminID")

		var products []model.Product
		if err := db.Preload("Media").Where("admin_id = ?", adminID).Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Could not fetch products",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"products": products,
			"count":    len(products),
		})
	}
}

// update products
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

// delete product
func DeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.GetUint("adminID")
		productID := c.Param("id")

		var product model.Product
		if err := db.Where("id = ? AND admin_id = ?", productID, adminID).First(&product).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not Found"})
			return
		}

		db.Delete(&product)
		c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
	}
}

// list all products - for user
func ListAllProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []model.Product
		if err := db.Preload("Media").Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}
