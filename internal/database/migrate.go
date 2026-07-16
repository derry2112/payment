package database

import (
	"payment/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Tag{},
		&model.Product{},
		&model.ProductDetail{},
		&model.ProductVariant{},
		&model.ProductImage{},
	)
}
