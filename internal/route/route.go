package route

import (
	"net/http"

	"payment/internal/handler"

	"github.com/gin-gonic/gin"
)

func New(productHandler *handler.ProductHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")
	{
		products := api.Group("/products")
		{
			products.POST("", productHandler.Create)
			products.GET("", productHandler.FindAll)
			products.GET("/:id", productHandler.FindByID)
			products.PATCH("/:id", productHandler.Update)
			products.DELETE("/:id", productHandler.Delete)
		}
	}

	return router
}
