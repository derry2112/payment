package route

import (
	"payment/internal/handler"

	"github.com/gin-gonic/gin"
)

func registerProductRoutes(
	api *gin.RouterGroup,
	productHandler *handler.ProductHandler,
) {
	products := api.Group("/products")

	products.POST("", productHandler.Create)
	products.GET("", productHandler.FindAll)
	products.GET("/:id", productHandler.FindByID)
	products.PATCH("/:id", productHandler.Update)
	products.DELETE("/:id", productHandler.Delete)
}
