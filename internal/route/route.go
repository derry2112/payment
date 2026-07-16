package route

import (
	"net/http"

	"payment/internal/handler"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	ProductHandler *handler.ProductHandler
}

func New(dependencies Dependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")
	registerProductRoutes(api, dependencies.ProductHandler)

	return router
}
