package products

import (
	"encoding/json"
	"net/http"

	"go-crud-api/internal/http/middleware"
	"go-crud-api/pkg/web"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ProductHandler handles product-related requests.
type ProductHandler struct {
	service  *Service
	validate *validator.Validate
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(service *Service) *ProductHandler {
	return &ProductHandler{
		service:  service,
		validate: validator.New(),
	}
}

// CreateProductRequest is the request payload for creating a product.
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=120"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
}

// UpdateProductRequest is the request payload for updating a product.
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=120"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
}

// CreateProduct handles product creation.
// @Summary Create a new product
// @Description Create a new product with name, description, price, and stock
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Product creation data"
// @Success 201 {object} web.Response{data=Product} "Product created successfully"
// @Failure 400 {object} web.Response{error=web.ApiError} "Bad request or validation error"
// @Failure 401 {object} web.Response{error=web.ApiError} "Unauthorized"
// @Failure 500 {object} web.Response{error=web.ApiError} "Internal server error"
// @Router /v1/products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondWithError(w, "bad_request", "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		web.RespondWithError(w, "validation_error", err.Error(), http.StatusBadRequest)
		return
	}

	ownerID, ok := r.Context().Value(middleware.ContextKeyUserID).(uuid.UUID)
	if !ok {
		web.RespondWithError(w, "unauthorized", "User ID not found in context", http.StatusUnauthorized)
		return
	}

	product, err := h.service.Create(r.Context(), req.Name, req.Description, req.Price, req.Stock, ownerID)
	if err != nil {
		web.RespondWithError(w, "internal_error", "Could not create product", http.StatusInternalServerError)
		return
	}

	web.RespondWithJSON(w, http.StatusCreated, web.Response{Data: product})
}

// GetProductByID handles fetching a product by ID.
// @Summary Get product by ID
// @Description Get product details by its ID
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Param productID path string true "Product ID"
// @Success 200 {object} web.Response{data=Product} "Product details"
// @Failure 400 {object} web.Response{error=web.ApiError} "Invalid product ID format"
// @Failure 401 {object} web.Response{error=web.ApiError} "Unauthorized"
// @Failure 404 {object} web.Response{error=web.ApiError} "Product not found"
// @Failure 500 {object} web.Response{error=web.ApiError} "Internal server error"
// @Router /v1/products/{productID} [get]
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "productID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondWithError(w, "bad_request", "Invalid product ID format", http.StatusBadRequest)
		return
	}

	product, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		// TODO: Differentiate between not found and other errors
		web.RespondWithError(w, "not_found", "Product not found", http.StatusNotFound)
		return
	}

	web.RespondWithJSON(w, http.StatusOK, web.Response{Data: product})
}

// ListProducts handles fetching all products.
// @Summary Get all products
// @Description Get a list of all products
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Success 200 {object} web.Response{data=[]Product} "List of products"
// @Failure 500 {object} web.Response{error=web.ApiError} "Internal server error"
// @Router /v1/products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.List(r.Context())
	if err != nil {
		web.RespondWithError(w, "internal_error", "Could not fetch products", http.StatusInternalServerError)
		return
	}

	web.RespondWithJSON(w, http.StatusOK, web.Response{Data: products})
}

// UpdateProduct handles updating an existing product.
// @Summary Update an existing product
// @Description Update product details by its ID. Only owner or admin can update.
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param productID path string true "Product ID"
// @Param product body UpdateProductRequest true "Product update data"
// @Success 200 {object} web.Response{data=Product} "Product updated successfully"
// @Failure 400 {object} web.Response{error=web.ApiError} "Bad request or validation error"
// @Failure 401 {object} web.Response{error=web.ApiError} "Unauthorized"
// @Failure 403 {object} web.Response{error=web.ApiError} "Forbidden (not owner or admin)"
// @Failure 404 {object} web.Response{error=web.ApiError} "Product not found"
// @Failure 500 {object} web.Response{error=web.ApiError} "Internal server error"
// @Router /v1/products/{productID} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "productID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondWithError(w, "bad_request", "Invalid product ID format", http.StatusBadRequest)
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondWithError(w, "bad_request", "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		web.RespondWithError(w, "validation_error", err.Error(), http.StatusBadRequest)
		return
	}

	ownerID, ok := r.Context().Value(middleware.ContextKeyUserID).(uuid.UUID)
	if !ok {
		web.RespondWithError(w, "unauthorized", "User ID not found in context", http.StatusUnauthorized)
		return
	}

	userRole, ok := r.Context().Value(middleware.ContextKeyRole).(string)
	if !ok {
		web.RespondWithError(w, "unauthorized", "User role not found in context", http.StatusUnauthorized)
		return
	}

	// Check ownership or admin role
	product, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		web.RespondWithError(w, "not_found", "Product not found", http.StatusNotFound)
		return
	}

	if product.OwnerID != ownerID && userRole != "admin" {
		web.RespondWithError(w, "forbidden", "You do not have permission to update this product", http.StatusForbidden)
		return
	}

	updatedProduct, err := h.service.Update(r.Context(), id, req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		web.RespondWithError(w, "internal_error", "Could not update product", http.StatusInternalServerError)
		return
	}

	web.RespondWithJSON(w, http.StatusOK, web.Response{Data: updatedProduct})
}

// DeleteProduct handles deleting a product.
// @Summary Delete a product
// @Description Delete a product by its ID. Only owner or admin can delete.
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Param productID path string true "Product ID"
// @Success 204 "Product deleted successfully"
// @Failure 400 {object} web.Response{error=web.ApiError} "Invalid product ID format"
// @Failure 401 {object} web.Response{error=web.ApiError} "Unauthorized"
// @Failure 403 {object} web.Response{error=web.ApiError} "Forbidden (not owner or admin)"
// @Failure 404 {object} web.Response{error=web.ApiError} "Product not found"
// @Failure 500 {object} web.Response{error=web.ApiError} "Internal server error"
// @Router /v1/products/{productID} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "productID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondWithError(w, "bad_request", "Invalid product ID format", http.StatusBadRequest)
		return
	}

	ownerID, ok := r.Context().Value(middleware.ContextKeyUserID).(uuid.UUID)
	if !ok {
		web.RespondWithError(w, "unauthorized", "User ID not found in context", http.StatusUnauthorized)
		return
	}

	userRole, ok := r.Context().Value(middleware.ContextKeyRole).(string)
	if !ok {
		web.RespondWithError(w, "unauthorized", "User role not found in context", http.StatusUnauthorized)
		return
	}

	// Check ownership or admin role
	product, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		web.RespondWithError(w, "not_found", "Product not found", http.StatusNotFound)
		return
	}

	if product.OwnerID != ownerID && userRole != "admin" {
		web.RespondWithError(w, "forbidden", "You do not have permission to delete this product", http.StatusForbidden)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		web.RespondWithError(w, "internal_error", "Could not delete product", http.StatusInternalServerError)
		return
	}

	web.RespondWithJSON(w, http.StatusNoContent, nil) // No content for successful delete
}