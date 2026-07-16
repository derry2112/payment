package dto

import "time"

type CreateProductRequest struct {
	CategoryID  *uint                         `json:"category_id" binding:"omitempty,min=1"`
	Name        string                        `json:"name" binding:"required,min=3,max=150"`
	Description string                        `json:"description" binding:"max=1000"`
	Price       int64                         `json:"price" binding:"gte=0"`
	Stock       int                           `json:"stock" binding:"gte=0"`
	Detail      *CreateProductDetailRequest   `json:"detail"`
	Variants    []CreateProductVariantRequest `json:"variants" binding:"max=50,dive"`
	Images      []CreateProductImageRequest   `json:"images" binding:"max=20,dive"`
	Tags        []CreateTagRequest            `json:"tags" binding:"max=20,dive"`
}

type UpdateProductRequest struct {
	CategoryID  *uint   `json:"category_id" binding:"omitempty,min=1"`
	Name        *string `json:"name" binding:"omitempty,min=3,max=150"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	Price       *int64  `json:"price" binding:"omitempty,gte=0"`
	Stock       *int    `json:"stock" binding:"omitempty,gte=0"`
}

type CreateProductDetailRequest struct {
	WeightGrams int     `json:"weight_grams" binding:"gte=0"`
	LengthCM    int     `json:"length_cm" binding:"gte=0"`
	WidthCM     int     `json:"width_cm" binding:"gte=0"`
	HeightCM    int     `json:"height_cm" binding:"gte=0"`
	Brand       string  `json:"brand" binding:"max=100"`
	SKU         *string `json:"sku" binding:"omitempty,max=100"`
}

type CreateProductVariantRequest struct {
	Name  string `json:"name" binding:"required,max=100"`
	SKU   string `json:"sku" binding:"required,max=100"`
	Price int64  `json:"price" binding:"gte=0"`
	Stock int    `json:"stock" binding:"gte=0"`
}

type CreateProductImageRequest struct {
	URL       string `json:"url" binding:"required,url"`
	AltText   string `json:"alt_text" binding:"max=200"`
	IsPrimary bool   `json:"is_primary"`
	Position  int    `json:"position" binding:"gte=0"`
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required,max=80"`
	Slug string `json:"slug" binding:"required,max=100"`
}

type ProductResponse struct {
	ID          uint                     `json:"id"`
	CategoryID  *uint                    `json:"category_id"`
	Category    *CategoryResponse        `json:"category,omitempty"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Price       int64                    `json:"price"`
	Stock       int                      `json:"stock"`
	Detail      *ProductDetailResponse   `json:"detail,omitempty"`
	Variants    []ProductVariantResponse `json:"variants"`
	Images      []ProductImageResponse   `json:"images"`
	Tags        []TagResponse            `json:"tags"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Meta     PaginationMeta    `json:"meta"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ProductDetailResponse struct {
	ID          uint    `json:"id"`
	WeightGrams int     `json:"weight_grams"`
	LengthCM    int     `json:"length_cm"`
	WidthCM     int     `json:"width_cm"`
	HeightCM    int     `json:"height_cm"`
	Brand       string  `json:"brand"`
	SKU         *string `json:"sku"`
}

type ProductVariantResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	SKU   string `json:"sku"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}

type ProductImageResponse struct {
	ID        uint   `json:"id"`
	URL       string `json:"url"`
	AltText   string `json:"alt_text"`
	IsPrimary bool   `json:"is_primary"`
	Position  int    `json:"position"`
}

type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
