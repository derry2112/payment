package repository

import (
	"context"

	"payment/internal/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	FindAll(ctx context.Context, limit, offset int) ([]model.Product, int64, error)
	FindByID(ctx context.Context, id uint) (*model.Product, error)
	CategoryExists(ctx context.Context, id uint) (bool, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, product *model.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return err
	}

	return preloadProductRelations(r.db.WithContext(ctx)).
		First(product, product.ID).Error
}

func (r *productRepository) FindAll(
	ctx context.Context,
	limit, offset int,
) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := preloadProductRelations(query).
		Order("products.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) FindByID(
	ctx context.Context,
	id uint,
) (*model.Product, error) {
	var product model.Product
	if err := preloadProductRelations(r.db.WithContext(ctx)).
		First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) CategoryExists(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Category{}).
		Where("id = ?", id).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	if err := r.db.WithContext(ctx).
		Omit("Category", "Detail", "Variants", "Images", "Tags").
		Save(product).Error; err != nil {
		return err
	}

	return preloadProductRelations(r.db.WithContext(ctx)).
		First(product, product.ID).Error
}

func (r *productRepository) Delete(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Delete(product).Error
}

func preloadProductRelations(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Category").
		Preload("Detail").
		Preload("Variants", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC, created_at ASC")
		}).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		})
}
