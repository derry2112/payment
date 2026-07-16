package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint             `gorm:"primaryKey"`
	CategoryID  *uint            `gorm:"index"`
	Category    *Category        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name        string           `gorm:"size:150;not null"`
	Description string           `gorm:"type:text"`
	Price       int64            `gorm:"not null;default:0"`
	Stock       int              `gorm:"not null;default:0"`
	Detail      *ProductDetail   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Variants    []ProductVariant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Images      []ProductImage   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags        []Tag            `gorm:"many2many:product_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
