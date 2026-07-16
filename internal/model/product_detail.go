package model

import "time"

type ProductDetail struct {
	ID          uint    `gorm:"primaryKey"`
	ProductID   uint    `gorm:"uniqueIndex;not null"`
	WeightGrams int     `gorm:"not null;default:0"`
	LengthCM    int     `gorm:"not null;default:0"`
	WidthCM     int     `gorm:"not null;default:0"`
	HeightCM    int     `gorm:"not null;default:0"`
	Brand       string  `gorm:"size:100"`
	SKU         *string `gorm:"size:100;uniqueIndex"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
