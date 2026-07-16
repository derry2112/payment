package mapper

import (
	"payment/internal/dto"
	"payment/internal/model"
)

func ToCategoryDetailResponse(category model.Category) dto.CategoryDetailResponse {
	return dto.CategoryDetailResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func ToCategoryDetailResponses(categories []model.Category) []dto.CategoryDetailResponse {
	responses := make([]dto.CategoryDetailResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, ToCategoryDetailResponse(category))
	}

	return responses
}
