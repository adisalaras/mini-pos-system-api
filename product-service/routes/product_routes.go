package routes

import (
	"product-service/handlers"
	"product-service/repositories"
	"product-service/services"

	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(app *fiber.App) {
	// inisialisasi layer/dependency
	productRepo := repositories.NewProductRepository()
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	api := app.Group("/api")
	products := api.Group("/products")
	
	products.Post("/", productHandler.CreateProduct)
	products.Get("/", productHandler.GetAllProducts)
	products.Get("/:id", productHandler.GetProduct)
	products.Put("/:id", productHandler.UpdateProduct)
	products.Delete("/:id", productHandler.DeleteProduct)
}
