package handler

import (
	"net/http"

	"github.com/arvtia/rest-api/internal/model"
	"github.com/arvtia/rest-api/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input model.Admin
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), bcrypt.DefaultCost)
		input.PasswordHash = string(hashed)

		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create admin"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Admin created"})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var admin model.Admin
		if err := db.Where("email = ?", input.Email).First(&admin).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, _ := utils.GenerateJWT(admin.ID, admin.Email)
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
