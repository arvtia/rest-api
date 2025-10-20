// handler/media.go
package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/arvtia/rest-api/internal/model"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadProductImage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")

		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
			return
		}

		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image"})
			return
		}
		defer src.Close()

		cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
		if err != nil {
			log.Println("Cloudinary config error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary config failed"})
			return
		}

		uploadResp, err := cld.Upload.Upload(context.Background(), src, uploader.UploadParams{
			Folder:   "products",
			PublicID: "product_" + productID,
		})
		if err != nil {
			log.Println("Upload error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
			return
		}

		media := model.ProductMedia{
			ProductID: parseUint(productID),
			URL:       uploadResp.SecureURL,
			Type:      "image",
			AltText:   file.Filename,
		}
		if err := db.Create(&media).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save media"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"url": media.URL})
	}
}

func parseUint(s string) uint {
	// safe conversion helper
	id, _ := strconv.ParseUint(s, 10, 64)
	return uint(id)
}
