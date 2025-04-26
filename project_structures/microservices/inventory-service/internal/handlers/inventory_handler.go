package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/edgardnogueira/go-patterns/project_structures/microservices/inventory-service/internal/domain"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/inventory-service/internal/repositories"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/common"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/messaging"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/observability"
	"github.com/go-chi/chi/v5"
)

// InventoryHandler handles HTTP requests for inventory
type InventoryHandler struct {
	repo      repositories.InventoryRepository
	publisher messaging.Publisher
	logger    *observability.Logger
}

// CreateProductRequest represents the request to create a new product
type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

// ReserveStockRequest represents the request to reserve stock
type ReserveStockRequest struct {
	OrderID  string                   `json:"order_id"`
	Items    []ReserveStockItemRequest `json:"items"`
}

// ReserveStockItemRequest represents an item in the stock reservation request
type ReserveStockItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(repo repositories.InventoryRepository, publisher messaging.Publisher, logger *observability.Logger) *InventoryHandler {
	return &InventoryHandler{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}

// ListProducts handles GET /api/inventory
func (h *InventoryHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := chi.URLParam(r, "request_id")
	if requestID == "" {
		requestID = common.GenerateRequestID()
	}

	// Get pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get products
	products, err := h.repo.ListProducts(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list products", err, map[string]interface{}{
			"request_id": requestID,
		})
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to list products", common.ErrInternalServerError, requestID)
		return
	}

	common.RespondWithJSON(w, http.StatusOK, products, requestID)
}

// GetProduct handles GET /api/inventory/{id}
func (h *InventoryHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := chi.URLParam(r, "request_id")
	if requestID == "" {
		requestID = common.GenerateRequestID()
	}

	// Get product ID from URL
	productID := chi.URLParam(r, "id")
	if productID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Product ID is required", common.ErrBadRequest, requestID)
		return
	}

	// Check if the client is checking availability
	if requiredStr := r.URL.Query().Get("required"); requiredStr != "" {
		required, err := strconv.Atoi(requiredStr)
		if err != nil || required <= 0 {
			common.RespondWithError(w, http.StatusBadRequest, "Invalid required quantity", common.ErrBadRequest, requestID)
			return
		}

		// Get product
		product, err := h.repo.GetProductByID(ctx, productID)
		if err != nil {
			h.logger.Error("Failed to get product", err, map[string]interface{}{
				"request_id": requestID,
				"product_id": productID,
			})
			common.RespondWithError(w, http.StatusNotFound, "Product not found", common.ErrNotFound, requestID)
			return
		}

		// Check availability
		if !product.IsAvailable(required) {
			common.RespondWithError(w, http.StatusConflict, "Insufficient inventory", common.ErrConflict, requestID)
			return
		}

		// Return success response
		common.RespondWithMessage(w, http.StatusOK, "Product is available", requestID)
		return
	}

	// Get product details
	product, err := h.repo.GetProductByID(ctx, productID)
	if err != nil {
		h.logger.Error("Failed to get product", err, map[string]interface{}{
			"request_id": requestID,
			"product_id": productID,
		})
		common.RespondWithError(w, http.StatusNotFound, "Product not found", common.ErrNotFound, requestID)
		return
	}

	common.RespondWithJSON(w, http.StatusOK, product, requestID)
}

// CreateProduct handles POST /api/inventory
func (h *InventoryHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := chi.URLParam(r, "request_id")
	if requestID == "" {
		requestID = common.GenerateRequestID()
	}

	// Parse request
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request payload", err, map[string]interface{}{
			"request_id": requestID,
		})
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", common.ErrBadRequest, requestID)
		return
	}

	// Create product
	product, err := domain.NewProduct(req.Name, req.Description, req.SKU, req.Price, req.Quantity)
	if err != nil {
		h.logger.Error("Failed to create product", err, map[string]interface{}{
			"request_id": requestID,
		})
		common.RespondWithError(w, http.StatusBadRequest, err.Error(), common.ErrBadRequest, requestID)
		return
	}

	// Save product to database
	if err := h.repo.CreateProduct(ctx, product); err != nil {
		h.logger.Error("Failed to save product", err, map[string]interface{}{
			"request_id": requestID,
		})
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to save product", common.ErrInternalServerError, requestID)
		return
	}

	common.RespondWithJSON(w, http.StatusCreated, product, requestID)
}

// UpdateProduct handles PUT /api/inventory/{id}
func (h *InventoryHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := chi.URLParam(r, "request_id")
	if requestID == "" {
		requestID = common.GenerateRequestID()
	}

	// Get product ID from URL
	productID := chi.URLParam(r, "id")
	if productID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Product ID is required", common.ErrBadRequest, requestID)
		return
	}

	// Parse request
	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request payload", err, map[string]interface{}{
			"request_id": requestID,
		})
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", common.ErrBadRequest, requestID)
		return
	}

	// Get existing product
	product, err := h.repo.GetProductByID(ctx, productID)
	if err != nil {
		h.logger.Error("Failed to get product", err, map[string]interface{}{
			"request_id": requestID,
			"product_id": productID,
		})
		common.RespondWithError(w, http.StatusNotFound, "Product not found", common.ErrNotFound, requestID)
		return
	}

	// Update product
	product.Name = req.Name
	product.Description = req.Description
	product.SKU = req.SKU
	product.Price = req.Price
	if err := product.UpdateStock(req.Quantity); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error(), common.ErrBadRequest, requestID)
		return
	}

	// Save updated product
	if err := h.repo.UpdateProduct(ctx, product); err != nil {
		h.logger.Error("Failed to update product", err, map[string]interface{}{
			"request_id": requestID,
			"product_id": productID,
		})
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to update product", common.ErrInternalServerError, requestID)
		return
	}

	common.RespondWithJSON(w, http.StatusOK, product, requestID)
}

// DeleteProduct handles DELETE /api/inventory/{id}
func (h *InventoryHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := chi.URLParam(r, "request_id")
	if requestID == "" {
		requestID = common.GenerateRequestID()
	}

	// Get product ID from URL
	productID := chi.URLParam(r, "id")
	if productID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Product ID is required", common.ErrBadRequest, requestID)
		return
	}

	// Delete product
	if err := h.repo.DeleteProduct(ctx, productID); err != nil {
		h.logger.Error("Failed to delete product", err, map[string]interface{}{
			"request_id": requestID,
			"product_id": productID,
		})
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to delete product", common.ErrInternalServerError, requestID)
		return
	}

	common.RespondWithMessage(w, http.StatusOK, "Product deleted successfully", requestID)
}

// ReserveStock handles POST /api/inventory/reserve
func (h *InventoryHandler) ReserveStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := chi.URLParam(r, "request_id")
	if requestID == "" {
		requestID = common.GenerateRequestID()
	}

	// Parse request
	var req ReserveStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request payload", err, map[string]interface{}{
			"request_id": requestID,
		})
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", common.ErrBadRequest, requestID)
		return
	}

	// Validate request
	if req.OrderID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Order ID is required", common.ErrBadRequest, requestID)
		return
	}

	if len(req.Items) == 0 {
		common.RespondWithError(w, http.StatusBadRequest, "At least one item is required", common.ErrBadRequest, requestID)
		return
	}

	// Process each item
	reservations := make([]*domain.InventoryReservation, 0, len(req.Items))
	for _, item := range req.Items {
		// Get product
		product, err := h.repo.GetProductByID(ctx, item.ProductID)
		if err != nil {
			h.logger.Error("Failed to get product", err, map[string]interface{}{
				"request_id": requestID,
				"product_id": item.ProductID,
			})
			common.RespondWithError(w, http.StatusNotFound, "Product not found", common.ErrNotFound, requestID)
			return
		}

		// Reserve stock
		reservation, err := product.ReserveStock(req.OrderID, item.Quantity)
		if err != nil {
			h.logger.Error("Failed to reserve stock", err, map[string]interface{}{
				"request_id": requestID,
				"product_id": item.ProductID,
				"order_id":   req.OrderID,
			})
			common.RespondWithError(w, http.StatusConflict, err.Error(), common.ErrConflict, requestID)
			return
		}

		// Save reservation
		if err := h.repo.CreateReservation(ctx, reservation); err != nil {
			h.logger.Error("Failed to save reservation", err, map[string]interface{}{
				"request_id": requestID,
				"product_id": item.ProductID,
				"order_id":   req.OrderID,
			})
			common.RespondWithError(w, http.StatusInternalServerError, "Failed to save reservation", common.ErrInternalServerError, requestID)
			return
		}

		// Update product stock
		if err := h.repo.UpdateProduct(ctx, product); err != nil {
			h.logger.Error("Failed to update product stock", err, map[string]interface{}{
				"request_id": requestID,
				"product_id": item.ProductID,
			})
			common.RespondWithError(w, http.StatusInternalServerError, "Failed to update product stock", common.ErrInternalServerError, requestID)
			return
		}

		reservations = append(reservations, reservation)
	}

	// Publish inventory reserved event
	err := h.publisher.Publish(messaging.InventoryReserved, map[string]interface{}{
		"order_id":     req.OrderID,
		"reservations": reservations,
	})
	if err != nil {
		h.logger.Error("Failed to publish inventory reserved event", err, map[string]interface{}{
			"request_id": requestID,
			"order_id":   req.OrderID,
		})
		// Continue even if event publishing fails
	}

	common.RespondWithJSON(w, http.StatusOK, reservations, requestID)
}
