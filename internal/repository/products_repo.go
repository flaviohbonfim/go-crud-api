package repository

import (
	"context"
	"go-crud-api/internal/domain/products"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormProductRepository struct {
	db *gorm.DB
}

// NewGormProductRepository creates a new GORM product repository.
func NewGormProductRepository(db *gorm.DB) products.ProductRepository {
	return &gormProductRepository{db: db}
}

func (r *gormProductRepository) Create(ctx context.Context, product *products.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *gormProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*products.Product, error) {
	var product products.Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *gormProductRepository) Update(ctx context.Context, product *products.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *gormProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&products.Product{}, id).Error
}
