package route

import (
	"github.com/gin-gonic/gin"

	"payment/internal/handler"
)

type Dependencies struct {
	ProductHandler  *handler.ProductHandler
	CategoryHandler *handler.CategoryHandler
	HealthHandler   *handler.HealthHandler
}

func New(dependencies Dependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	registerHealthRoutes(router, dependencies.HealthHandler)

	api := router.Group("/api")
	registerProductRoutes(api, dependencies.ProductHandler)
	registerCategoryRoutes(api, dependencies.CategoryHandler)

	return router
}
