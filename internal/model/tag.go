package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:80;not null"`
	Slug      string    `gorm:"size:100;uniqueIndex;not null"`
	Products  []Product `gorm:"many2many:product_tags;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
