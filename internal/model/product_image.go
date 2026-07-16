package model

import "time"

type ProductImage struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   `gorm:"index;not null"`
	URL       string `gorm:"type:text;not null"`
	AltText   string `gorm:"size:200"`
	IsPrimary bool   `gorm:"not null;default:false"`
	Position  int    `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
