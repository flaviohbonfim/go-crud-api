package http

import (
	"go-crud-api/internal/config"
	"go-crud-api/internal/domain/products"
	"go-crud-api/internal/domain/users"
	"go-crud-api/internal/http/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2" // <- alias necessário
	"gorm.io/gorm"

	_ "go-crud-api/docs"
)

// InitRouter initializes and returns a new chi router.
func InitRouter(cfg config.Config, db *gorm.DB, authHandler *users.AuthHandler, productHandler *products.ProductHandler) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	// r.Use(chimiddleware.Logger) // We'll use our custom zerolog logger
	r.Use(chimiddleware.Recoverer)

	// Health check endpoint
	r.Get("/healthz", HealthCheckHandler(db))

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	// Auth routes
	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg))

		// User routes (Admin only)
		r.Route("/v1/users", func(r chi.Router) {
			r.Use(middleware.HasRoleMiddleware("admin"))
			r.Get("/", authHandler.ListUsers)
		})

		// Product routes
		r.Route("/v1/products", func(r chi.Router) {
			r.Get("/", productHandler.ListProducts)
			r.Post("/", productHandler.CreateProduct)
			r.Get("/{productID}", productHandler.GetProductByID)
			r.Put("/{productID}", productHandler.UpdateProduct)
			r.Delete("/{productID}", productHandler.DeleteProduct)
		})
	})

	return r
}
