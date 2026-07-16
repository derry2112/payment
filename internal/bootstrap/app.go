package bootstrap

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"payment/internal/handler"
	"payment/internal/repository"
	"payment/internal/route"
	"payment/internal/service"
)

type App struct {
	Router *gin.Engine
}

func NewApp(db *gorm.DB) *App {
	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	categoryRepository := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	healthHandler := handler.NewHealthHandler(db)

	router := route.New(route.Dependencies{
		ProductHandler:  productHandler,
		CategoryHandler: categoryHandler,
		HealthHandler:   healthHandler,
	})

	return &App{
		Router: router,
	}
}
