package model

import "time"

type Order struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	TotalAmount int64
	Currency    string
	Status      string // paid pending failed
	PaymentID   string // razerpay / stripe ID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`
}
