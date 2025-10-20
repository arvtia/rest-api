package model

import (
	"time"
)

type ProductMedia struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   `gorm:"index"`
	URL       string `json:"url"`
	Type      string `json:"type"`
	AltText   string `json:"alt_text"`
	CreatedAt time.Time
}
