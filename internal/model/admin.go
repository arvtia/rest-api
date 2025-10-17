package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"uniqueIndex"`
	PasswordHash string `json:"-"`
	StoreName    string `json:"store_name"`
}


