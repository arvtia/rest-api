package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Stock       int            `json:"stock"`
	Category    string         `json:"category"`
	AdminID     uint           `json:"admin_id"`
	Media       []ProductMedia `gorm:"foreignKey:ProductID"`
}
