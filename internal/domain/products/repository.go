package products

import (
	"context"

	"github.com/google/uuid"
)

// ProductRepository defines the interface for product data operations.
 type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]Product, error)
	// TODO: Add List method with filters and pagination
}
