package mapper

import (
	"payment/internal/dto"
	"payment/internal/model"
)

func ToProductResponse(product model.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Category:    toCategoryResponse(product.Category),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Detail:      toProductDetailResponse(product.Detail),
		Variants:    toProductVariantResponses(product.Variants),
		Images:      toProductImageResponses(product.Images),
		Tags:        toTagResponses(product.Tags),
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func ToProductResponses(products []model.Product) []dto.ProductResponse {
	responses := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, ToProductResponse(product))
	}

	return responses
}

func toCategoryResponse(category *model.Category) *dto.CategoryResponse {
	if category == nil {
		return nil
	}

	return &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}
}

func toProductDetailResponse(detail *model.ProductDetail) *dto.ProductDetailResponse {
	if detail == nil {
		return nil
	}

	return &dto.ProductDetailResponse{
		ID:          detail.ID,
		WeightGrams: detail.WeightGrams,
		LengthCM:    detail.LengthCM,
		WidthCM:     detail.WidthCM,
		HeightCM:    detail.HeightCM,
		Brand:       detail.Brand,
		SKU:         detail.SKU,
	}
}

func toProductVariantResponses(variants []model.ProductVariant) []dto.ProductVariantResponse {
	responses := make([]dto.ProductVariantResponse, 0, len(variants))
	for _, variant := range variants {
		responses = append(responses, dto.ProductVariantResponse{
			ID:    variant.ID,
			Name:  variant.Name,
			SKU:   variant.SKU,
			Price: variant.Price,
			Stock: variant.Stock,
		})
	}

	return responses
}

func toProductImageResponses(images []model.ProductImage) []dto.ProductImageResponse {
	responses := make([]dto.ProductImageResponse, 0, len(images))
	for _, image := range images {
		responses = append(responses, dto.ProductImageResponse{
			ID:        image.ID,
			URL:       image.URL,
			AltText:   image.AltText,
			IsPrimary: image.IsPrimary,
			Position:  image.Position,
		})
	}

	return responses
}

func toTagResponses(tags []model.Tag) []dto.TagResponse {
	responses := make([]dto.TagResponse, 0, len(tags))
	for _, tag := range tags {
		responses = append(responses, dto.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	return responses
}
