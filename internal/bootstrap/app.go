package bootstrap

import (
	"payment/internal/handler"
	"payment/internal/repository"
	"payment/internal/route"
	"payment/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
}

func NewApp(db *gorm.DB) *App {
	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	router := route.New(route.Dependencies{
		ProductHandler: productHandler,
	})

	return &App{
		Router: router,
	}
}
