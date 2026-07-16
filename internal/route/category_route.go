package route

import (
	"payment/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerCategoryRoutes(
	api *gin.RouterGroup,
	categoryHandler *handler.CategoryHandler,
) {
	categories := api.Group("/category")

	categories.POST("", categoryHandler.Create)
	categories.GET("", categoryHandler.FindAll)
	categories.GET("/:id", categoryHandler.FindByID)
	categories.GET("/:id/products", categoryHandler.FindProducts)
	categories.PATCH("/:id", categoryHandler.Update)
	categories.DELETE("/:id", categoryHandler.Delete)
}
