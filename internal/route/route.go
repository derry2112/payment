package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"payment/internal/handler"
)

type Dependencies struct {
	ProductHandler  *handler.ProductHandler
	CategoryHandler *handler.CategoryHandler
}

func New(dependencies Dependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api")
	registerProductRoutes(api, dependencies.ProductHandler)
	registerCategoryRoutes(api, dependencies.CategoryHandler)

	return router
}
