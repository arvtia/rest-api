package config

import (
	"log"

	"github.com/arvtia/rest-api/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dialector := sqlite.Open("file:db.sqlite?mode=rwc&_pragma=foreign_keys(1)")
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&model.Admin{}, &model.Product{})
	if err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	return db
}
