package handler

import (
	"errors"
	"net/http"
	"strconv"

	"payment/internal/dto"
	"payment/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var request dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeValidationError(c, err)
		return
	}

	category, err := h.service.Create(c.Request.Context(), request)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	writeSuccess(c, http.StatusCreated, "kategori berhasil dibuat", category)
}

func (h *CategoryHandler) FindAll(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		writeValidationError(c, err)
		return
	}

	categories, err := h.service.FindAll(c.Request.Context(), pagination)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, "daftar kategori berhasil diambil", categories)
}

func (h *CategoryHandler) FindByID(c *gin.Context) {
	id, ok := categoryID(c)
	if !ok {
		return
	}

	category, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, "kategori berhasil diambil", category)
}

func (h *CategoryHandler) FindProducts(c *gin.Context) {
	id, ok := categoryID(c)
	if !ok {
		return
	}

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		writeValidationError(c, err)
		return
	}

	products, err := h.service.FindProducts(c.Request.Context(), id, pagination)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, "produk kategori berhasil diambil", products)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, ok := categoryID(c)
	if !ok {
		return
	}

	var request dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeValidationError(c, err)
		return
	}

	category, err := h.service.Update(c.Request.Context(), id, request)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, "kategori berhasil diperbarui", category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, ok := categoryID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.writeServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CategoryHandler) writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrCategoryNotFound):
		writeError(c, http.StatusNotFound, err.Error(), nil)
	case errors.Is(err, service.ErrNoCategoryChanges):
		writeError(c, http.StatusBadRequest, err.Error(), nil)
	default:
		writeError(c, http.StatusInternalServerError, "terjadi kesalahan pada server", nil)
	}
}

func categoryID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		writeError(c, http.StatusBadRequest, "ID kategori tidak valid", nil)
		return 0, false
	}

	return uint(id), true
}
