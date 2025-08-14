package products

import (
	"context"

	"github.com/google/uuid"
)

// Service defines the product service.
type Service struct {
	repo ProductRepository
}

// NewService creates a new product service.
func NewService(repo ProductRepository) *Service {
	return &Service{repo: repo}
}

// Create creates a new product.
func (s *Service) Create(ctx context.Context, name, description string, price float64, stock int, ownerID uuid.UUID) (*Product, error) {
	product := &Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		OwnerID:     ownerID,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// FindByID finds a product by its ID.
func (s *Service) FindByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	return s.repo.FindByID(ctx, id)
}

// Update updates a product.
func (s *Service) Update(ctx context.Context, id uuid.UUID, name, description string, price float64, stock int) (*Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	product.Name = name
	product.Description = description
	product.Price = price
	product.Stock = stock

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// Delete deletes a product.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
