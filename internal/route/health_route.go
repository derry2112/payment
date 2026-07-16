package route

import (
	"github.com/gin-gonic/gin"

	"payment/internal/handler"
)

func registerHealthRoutes(router *gin.Engine, healthHandler *handler.HealthHandler) {
	// router.GET("/health", healthHandler.Live)
	router.GET("/health/live", healthHandler.Live)
	router.GET("/health/ready", healthHandler.Ready)
}
