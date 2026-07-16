package service

import (
	"context"
	"errors"

	"payment/internal/dto"
	"payment/internal/mapper"
	"payment/internal/model"
	"payment/internal/repository"

	"gorm.io/gorm"
)

type CategoryService interface {
	Create(ctx context.Context, request dto.CreateCategoryRequest) (*dto.CategoryDetailResponse, error)
	FindAll(ctx context.Context, pagination dto.PaginationRequest) (*dto.CategoryListResponse, error)
	FindByID(ctx context.Context, id uint) (*dto.CategoryDetailResponse, error)
	FindProducts(ctx context.Context, id uint, pagination dto.PaginationRequest) (*dto.CategoryProductsResponse, error)
	Update(ctx context.Context, id uint, request dto.UpdateCategoryRequest) (*dto.CategoryDetailResponse, error)
	Delete(ctx context.Context, id uint) error
}

type categoryService struct {
	repository repository.CategoryRepository
}

func NewCategoryService(repository repository.CategoryRepository) CategoryService {
	return &categoryService{repository: repository}
}

func (s *categoryService) Create(
	ctx context.Context,
	request dto.CreateCategoryRequest,
) (*dto.CategoryDetailResponse, error) {
	category := model.Category{Name: request.Name, Slug: request.Slug}
	if err := s.repository.Create(ctx, &category); err != nil {
		return nil, err
	}

	response := mapper.ToCategoryDetailResponse(category)
	return &response, nil
}

func (s *categoryService) FindAll(
	ctx context.Context,
	pagination dto.PaginationRequest,
) (*dto.CategoryListResponse, error) {
	categories, total, err := s.repository.FindAll(
		ctx,
		pagination.Limit,
		pagination.Offset(),
	)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryListResponse{
		Categories: mapper.ToCategoryDetailResponses(categories),
		Meta:       dto.NewPaginationMeta(pagination, total),
	}, nil
}

func (s *categoryService) FindByID(
	ctx context.Context,
	id uint,
) (*dto.CategoryDetailResponse, error) {
	category, err := s.findCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mapper.ToCategoryDetailResponse(*category)
	return &response, nil
}

func (s *categoryService) FindProducts(
	ctx context.Context,
	id uint,
	pagination dto.PaginationRequest,
) (*dto.CategoryProductsResponse, error) {
	category, err := s.findCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	products, total, err := s.repository.FindProducts(
		ctx,
		id,
		pagination.Limit,
		pagination.Offset(),
	)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryProductsResponse{
		Category: mapper.ToCategoryDetailResponse(*category),
		Products: mapper.ToProductResponses(products),
		Meta:     dto.NewPaginationMeta(pagination, total),
	}, nil
}

func (s *categoryService) Update(
	ctx context.Context,
	id uint,
	request dto.UpdateCategoryRequest,
) (*dto.CategoryDetailResponse, error) {
	if request.Name == nil && request.Slug == nil {
		return nil, ErrNoCategoryChanges
	}

	category, err := s.findCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	if request.Name != nil {
		category.Name = *request.Name
	}
	if request.Slug != nil {
		category.Slug = *request.Slug
	}

	if err := s.repository.Update(ctx, category); err != nil {
		return nil, err
	}

	response := mapper.ToCategoryDetailResponse(*category)
	return &response, nil
}

func (s *categoryService) Delete(ctx context.Context, id uint) error {
	category, err := s.findCategory(ctx, id)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, category)
}

func (s *categoryService) findCategory(
	ctx context.Context,
	id uint,
) (*model.Category, error) {
	category, err := s.repository.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCategoryNotFound
	}
	if err != nil {
		return nil, err
	}

	return category, nil
}
