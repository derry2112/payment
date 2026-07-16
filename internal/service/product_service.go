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

type ProductService interface {
	Create(ctx context.Context, request dto.CreateProductRequest) (*dto.ProductResponse, error)
	FindAll(ctx context.Context, pagination dto.PaginationRequest) (*dto.ProductListResponse, error)
	FindByID(ctx context.Context, id uint) (*dto.ProductResponse, error)
	Update(ctx context.Context, id uint, request dto.UpdateProductRequest) (*dto.ProductResponse, error)
	Delete(ctx context.Context, id uint) error
}

type productService struct {
	repository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) ProductService {
	return &productService{repository: repository}
}

func (s *productService) Create(
	ctx context.Context,
	request dto.CreateProductRequest,
) (*dto.ProductResponse, error) {
	if err := s.validateCategory(ctx, request.CategoryID); err != nil {
		return nil, err
	}

	product := model.Product{
		CategoryID:  request.CategoryID,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
		Detail:      toProductDetail(request.Detail),
		Variants:    toProductVariants(request.Variants),
		Images:      toProductImages(request.Images),
		Tags:        toTags(request.Tags),
	}

	if err := s.repository.Create(ctx, &product); err != nil {
		return nil, err
	}

	response := mapper.ToProductResponse(product)
	return &response, nil
}

func (s *productService) FindAll(
	ctx context.Context,
	pagination dto.PaginationRequest,
) (*dto.ProductListResponse, error) {
	products, total, err := s.repository.FindAll(
		ctx,
		pagination.Limit,
		pagination.Offset(),
	)
	if err != nil {
		return nil, err
	}

	return &dto.ProductListResponse{
		Products: mapper.ToProductResponses(products),
		Meta:     dto.NewPaginationMeta(pagination, total),
	}, nil
}

func (s *productService) FindByID(
	ctx context.Context,
	id uint,
) (*dto.ProductResponse, error) {
	product, err := s.findProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	response := mapper.ToProductResponse(*product)
	return &response, nil
}

func (s *productService) Update(
	ctx context.Context,
	id uint,
	request dto.UpdateProductRequest,
) (*dto.ProductResponse, error) {
	if request.Name == nil &&
		request.CategoryID == nil &&
		request.Description == nil &&
		request.Price == nil &&
		request.Stock == nil {
		return nil, ErrNoProductChanges
	}

	if err := s.validateCategory(ctx, request.CategoryID); err != nil {
		return nil, err
	}

	product, err := s.findProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	if request.Name != nil {
		product.Name = *request.Name
	}
	if request.CategoryID != nil {
		product.CategoryID = request.CategoryID
	}
	if request.Description != nil {
		product.Description = *request.Description
	}
	if request.Price != nil {
		product.Price = *request.Price
	}
	if request.Stock != nil {
		product.Stock = *request.Stock
	}

	if err := s.repository.Update(ctx, product); err != nil {
		return nil, err
	}

	response := mapper.ToProductResponse(*product)
	return &response, nil
}

func (s *productService) Delete(ctx context.Context, id uint) error {
	product, err := s.findProduct(ctx, id)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, product)
}

func (s *productService) findProduct(
	ctx context.Context,
	id uint,
) (*model.Product, error) {
	product, err := s.repository.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) validateCategory(ctx context.Context, categoryID *uint) error {
	if categoryID == nil {
		return nil
	}

	exists, err := s.repository.CategoryExists(ctx, *categoryID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCategoryNotFound
	}

	return nil
}

func toProductDetail(request *dto.CreateProductDetailRequest) *model.ProductDetail {
	if request == nil {
		return nil
	}

	return &model.ProductDetail{
		WeightGrams: request.WeightGrams,
		LengthCM:    request.LengthCM,
		WidthCM:     request.WidthCM,
		HeightCM:    request.HeightCM,
		Brand:       request.Brand,
		SKU:         request.SKU,
	}
}

func toProductVariants(requests []dto.CreateProductVariantRequest) []model.ProductVariant {
	variants := make([]model.ProductVariant, 0, len(requests))
	for _, request := range requests {
		variants = append(variants, model.ProductVariant{
			Name:  request.Name,
			SKU:   request.SKU,
			Price: request.Price,
			Stock: request.Stock,
		})
	}

	return variants
}

func toProductImages(requests []dto.CreateProductImageRequest) []model.ProductImage {
	images := make([]model.ProductImage, 0, len(requests))
	for _, request := range requests {
		images = append(images, model.ProductImage{
			URL:       request.URL,
			AltText:   request.AltText,
			IsPrimary: request.IsPrimary,
			Position:  request.Position,
		})
	}

	return images
}

func toTags(requests []dto.CreateTagRequest) []model.Tag {
	tags := make([]model.Tag, 0, len(requests))
	for _, request := range requests {
		tags = append(tags, model.Tag{
			Name: request.Name,
			Slug: request.Slug,
		})
	}

	return tags
}
