package handler

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"payment/internal/dto"
	"payment/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var request dto.CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeValidationError(c, err)
		return
	}

	product, err := h.service.Create(c.Request.Context(), request)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	writeSuccess(c, http.StatusCreated, "produk berhasil dibuat", product)
}

func (h *ProductHandler) FindAll(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		writeValidationError(c, err)
		return
	}

	products, err := h.service.FindAll(c.Request.Context(), pagination)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "gagal mengambil produk", nil)
		return
	}

	writeSuccess(c, http.StatusOK, "daftar produk berhasil diambil", products)
}

func (h *ProductHandler) FindByID(c *gin.Context) {
	id, ok := productID(c)
	if !ok {
		return
	}

	product, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, "produk berhasil diambil", product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, ok := productID(c)
	if !ok {
		return
	}

	var request dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeValidationError(c, err)
		return
	}

	product, err := h.service.Update(c.Request.Context(), id, request)
	if err != nil {
		h.writeServiceError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, "produk berhasil diperbarui", product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, ok := productID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.writeServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrProductNotFound):
		writeError(c, http.StatusNotFound, err.Error(), nil)
	case errors.Is(err, service.ErrNoProductChanges):
		writeError(c, http.StatusBadRequest, err.Error(), nil)
	case errors.Is(err, service.ErrCategoryNotFound):
		writeError(c, http.StatusUnprocessableEntity, err.Error(), nil)
	default:
		writeError(c, http.StatusInternalServerError, "terjadi kesalahan pada server", nil)
	}
}

func productID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		writeError(c, http.StatusBadRequest, "ID produk tidak valid", nil)
		return 0, false
	}

	return uint(id), true
}

func writeValidationError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		details := make(map[string]string, len(validationErrors))
		for _, fieldError := range validationErrors {
			details[fieldError.Field()] = validationMessage(fieldError)
		}
		writeError(c, http.StatusBadRequest, "data request tidak valid", details)
		return
	}

	writeError(c, http.StatusBadRequest, "format JSON tidak valid", nil)
}

func validationMessage(field validator.FieldError) string {
	switch field.Tag() {
	case "required":
		return "wajib diisi"
	case "min":
		if field.Kind() != reflect.String {
			return "minimal bernilai " + field.Param()
		}
		return "minimal " + field.Param() + " karakter"
	case "max":
		if field.Kind() != reflect.String {
			return "maksimal bernilai " + field.Param()
		}
		return "maksimal " + field.Param() + " karakter"
	case "gte":
		return "minimal bernilai " + field.Param()
	default:
		return "tidak valid"
	}
}
