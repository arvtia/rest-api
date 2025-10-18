package config

import (
	"log"

	"github.com/arvtia/rest-api/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// auto migrate models
	err = db.AutoMigrate(&model.Admin{}, &model.Product{})
	if err != nil {
		log.Fatal("Failed to migrate database", err)
	}
	return db
}
