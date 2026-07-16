package model

import "time"

type ProductVariant struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   `gorm:"index;not null"`
	Name      string `gorm:"size:100;not null"`
	SKU       string `gorm:"size:100;uniqueIndex;not null"`
	Price     int64  `gorm:"not null;default:0"`
	Stock     int    `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
