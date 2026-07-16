package repository

import (
	"context"

	"payment/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	FindAll(ctx context.Context, limit, offset int) ([]model.Category, int64, error)
	FindByID(ctx context.Context, id uint) (*model.Category, error)
	FindProducts(ctx context.Context, categoryID uint, limit, offset int) ([]model.Product, int64, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, category *model.Category) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) FindAll(
	ctx context.Context,
	limit, offset int,
) ([]model.Category, int64, error) {
	var categories []model.Category
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Category{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) FindByID(
	ctx context.Context,
	id uint,
) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) FindProducts(
	ctx context.Context,
	categoryID uint,
	limit, offset int,
) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("category_id = ?", categoryID)
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

func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Product{}).
			Where("category_id = ?", category.ID).
			Update("category_id", nil).Error; err != nil {
			return err
		}

		return tx.Delete(category).Error
	})
}
