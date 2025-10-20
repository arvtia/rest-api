package model

import (
	"time"
)

type CartItems struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	ProductID uint
	Quantity  int
	AddedAt   time.Time
}
