package products

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockProductRepository is a mock implementation of ProductRepository.
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestProductService_Create(t *testing.T) {
	repo := new(MockProductRepository)
	service := NewService(repo)

	ctx := context.Background()
	ownerID := uuid.New()

	// Test case 1: Successful creation
	repo.On("Create", ctx, mock.AnythingOfType("*products.Product")).Return(nil).Once()
	product, err := service.Create(ctx, "Test Product", "Desc", 10.50, 5, ownerID)
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Test Product", product.Name)
	repo.AssertExpectations(t)

	// Test case 2: Repository returns an error
	repo.On("Create", ctx, mock.AnythingOfType("*products.Product")).Return(errors.New("db error")).Once()
	product, err = service.Create(ctx, "Test Product", "Desc", 10.50, 5, ownerID)
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "db error")
	repo.AssertExpectations(t)
}

func TestProductService_FindByID(t *testing.T) {
	repo := new(MockProductRepository)
	service := NewService(repo)

	ctx := context.Background()
	productID := uuid.New()

	// Test case 1: Product found
	expectedProduct := &Product{ID: productID, Name: "Found Product"}
	repo.On("FindByID", ctx, productID).Return(expectedProduct, nil).Once()
	product, err := service.FindByID(ctx, productID)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	repo.AssertExpectations(t)

	// Test case 2: Product not found
	repo.On("FindByID", ctx, productID).Return(nil, gorm.ErrRecordNotFound).Once()
	product, err = service.FindByID(ctx, productID)
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), gorm.ErrRecordNotFound.Error())
	repo.AssertExpectations(t)
}

func TestProductService_Update(t *testing.T) {
	repo := new(MockProductRepository)
	service := NewService(repo)

	ctx := context.Background()
	productID := uuid.New()
	ownerID := uuid.New()

	// Test case 1: Successful update
	existingProduct := &Product{ID: productID, Name: "Old Name", OwnerID: ownerID}
	repo.On("FindByID", ctx, productID).Return(existingProduct, nil).Once()
	repo.On("Update", ctx, mock.AnythingOfType("*products.Product")).Return(nil).Once()
	updatedProduct, err := service.Update(ctx, productID, "New Name", "New Desc", 20.0, 10)
	assert.NoError(t, err)
	assert.NotNil(t, updatedProduct)
	assert.Equal(t, "New Name", updatedProduct.Name)
	repo.AssertExpectations(t)

	// Test case 2: Product not found
	repo.On("FindByID", ctx, productID).Return(nil, gorm.ErrRecordNotFound).Once()
	updatedProduct, err = service.Update(ctx, productID, "New Name", "New Desc", 20.0, 10)
	assert.Error(t, err)
	assert.Nil(t, updatedProduct)
	assert.Contains(t, err.Error(), gorm.ErrRecordNotFound.Error())
	repo.AssertExpectations(t)

	// Test case 3: Repository update error
	repo.On("FindByID", ctx, productID).Return(existingProduct, nil).Once()
	repo.On("Update", ctx, mock.AnythingOfType("*products.Product")).Return(errors.New("db error")).Once()
	updatedProduct, err = service.Update(ctx, productID, "New Name", "New Desc", 20.0, 10)
	assert.Error(t, err)
	assert.Nil(t, updatedProduct)
	assert.Contains(t, err.Error(), "db error")
	repo.AssertExpectations(t)
}

func TestProductService_Delete(t *testing.T) {
	repo := new(MockProductRepository)
	service := NewService(repo)

	ctx := context.Background()
	productID := uuid.New()

	// Test case 1: Successful delete
	repo.On("Delete", ctx, productID).Return(nil).Once()
	err := service.Delete(ctx, productID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)

	// Test case 2: Repository returns an error
	repo.On("Delete", ctx, productID).Return(errors.New("db error")).Once()
	err = service.Delete(ctx, productID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	repo.AssertExpectations(t)
}
