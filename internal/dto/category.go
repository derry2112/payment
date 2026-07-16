package dto

import "time"

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
	Slug string `json:"slug" binding:"required,min=3,max=120"`
}

type UpdateCategoryRequest struct {
	Name *string `json:"name" binding:"omitempty,min=3,max=100"`
	Slug *string `json:"slug" binding:"omitempty,min=3,max=120"`
}

type CategoryDetailResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryListResponse struct {
	Categories []CategoryDetailResponse `json:"categories"`
	Meta       PaginationMeta           `json:"meta"`
}

type CategoryProductsResponse struct {
	Category CategoryDetailResponse `json:"category"`
	Products []ProductResponse      `json:"products"`
	Meta     PaginationMeta         `json:"meta"`
}
