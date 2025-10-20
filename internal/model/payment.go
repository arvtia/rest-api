package model

import "time"

type Payment struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint
	OrderID uint

	Gateway       string // strip
	TransactionID string
	Status        string
	Amount        int64
	CreatedAt     time.Time
}
