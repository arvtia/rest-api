package model

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:100"`
	Email        string `gorm:"uniqueIndex;size:100"`
	PasswordHash string `gorm:"size:255"`
	Phone        string `gorm:"size:20"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserDetails  UserDetails `gorm:"foreignKey:UserID"`
	CartItems    []CartItems `gorm:"foreignKey:UserID"`
	Orders       []Order     `gorm:"foreignKey:UserID"`
}
