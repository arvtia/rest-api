package model

type UserDetails struct {
	ID         uint `gorm:"PrimaryKey"`
	UserID     uint `gorm:"uniqueIndex"` //one to one
	Address    string
	City       string
	State      string
	PostalCode string
	Country    string
	IsDefault  bool
}
