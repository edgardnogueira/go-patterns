package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/microservices/order-service/internal/domain"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/order-service/internal/repositories"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/common"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/messaging"
	"github.com/edgardnogueira/go-patterns/project_structures/microservices/pkg/observability"
	"github.com/go-chi/chi/v5"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	repo            repositories.OrderRepository
	publisher       messaging.Publisher
	inventoryClient *http.Client
	inventoryURL    string
	logger          *observability.Logger
}

// CreateOrderRequest represents the request to create a new order
type CreateOrderRequest struct {
	CustomerID string              `json:"customer_id"`
	Items      []CreateOrderItem   `json:"items"`
}

// CreateOrderItem represents an item in the order creation request
type CreateOrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// UpdateOrderRequest represents the request to update an order
type UpdateOrderRequest struct {
	Status string `json:"status"`
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(repo repositories.OrderRepository, publisher messaging.Publisher, inventoryURL string, logger *observability.Logger) *OrderHandler {
	return &OrderHandler{
		repo:            repo,
		publisher:       publisher,
		inventoryClient: &http.Client{Timeout: 5 * time.Second},
		inventoryURL:    inventoryURL,
		logger:          logger,
	}
}

// ListOrders handles GET /api/orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
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

	// Get orders
	orders, err := h.repo.ListOrders(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list orders", err, map[string]interface{}{
			"request_id": requestID,
		})
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to list orders", common.ErrInternalServerError, requestID)
		return
	}

	common.RespondWithJSON(w, http.StatusOK, orders, requestID)
}

// CreateOrder handles POST /api/orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := chi.URLParam(r, "request_id")
	if requestID == "" {
		requestID = common.GenerateRequestID()
	}

	// Parse request
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request payload", err, map[string]interface{}{
			"request_id": requestID,
		})
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", common.ErrBadRequest, requestID)
		return
	}

	// Validate request
	if req.CustomerID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Customer ID is required", common.ErrBadRequest, requestID)
		return
	}

	if len(req.Items) == 0 {
		common.RespondWithError(w, http.StatusBadRequest, "At least one item is required", common.ErrBadRequest, requestID)
		return
	}

	// Check inventory for each item
	for _, item := range req.Items {
		available, err := h.checkInventory(ctx, item.ProductID, item.Quantity)
		if err != nil {
			h.logger.Error("Failed to check inventory", err, map[string]interface{}{
				"request_id": requestID,
				"product_id": item.ProductID,
			})
			common.RespondWithError(w, http.StatusInternalServerError, "Failed to check inventory", common.ErrInternalServerError, requestID)
			return
		}

		if !available {
			common.RespondWithError(w, http.StatusConflict, fmt.Sprintf("Insufficient inventory for product %s", item.ProductID), common.ErrConflict, requestID)
			return
		}
	}

	// Convert request items to domain items
	orderItems := make([]domain.OrderItem, len(req.Items))
	for i, item := range req.Items {
		orderItems[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	// Create order
	order, err := domain.CreateOrder(req.CustomerID, orderItems)
	if err != nil {
		h.logger.Error("Failed to create order", err, map[string]interface{}{
			"request_id": requestID,
		})
		common.RespondWithError(w, http.StatusBadRequest, err.Error(), common.ErrBadRequest, requestID)
		return
	}

	// Save order to database
	if err := h.repo.CreateOrder(ctx, order); err != nil {
		h.logger.Error("Failed to save order", err, map[string]interface{}{
			"request_id": requestID,
		})
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to save order", common.ErrInternalServerError, requestID)
		return
	}

	// Publish order created event
	if err := h.publishOrderCreatedEvent(order); err != nil {
		h.logger.Error("Failed to publish order created event", err, map[string]interface{}{
			"request_id": requestID,
			"order_id":   order.ID,
		})
		// Continue even if the event publishing fails
	}

	common.RespondWithJSON(w, http.StatusCreated, order, requestID)
}

// GetOrder handles GET /api/orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := chi.URLParam(r, "request_id")
	if requestID == "" {
		requestID = common.GenerateRequestID()
	}

	// Get order ID from URL
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Order ID is required", common.ErrBadRequest, requestID)
		return
	}

	// Get order
	order, err := h.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		h.logger.Error("Failed to get order", err, map[string]interface{}{
			"request_id": requestID,
			"order_id":   orderID,
		})
		common.RespondWithError(w, http.StatusNotFound, "Order not found", common.ErrNotFound, requestID)
		return
	}

	common.RespondWithJSON(w, http.StatusOK, order, requestID)
}

// UpdateOrder handles PUT /api/orders/{id}
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := chi.URLParam(r, "request_id")
	if requestID == "" {
		requestID = common.GenerateRequestID()
	}

	// Get order ID from URL
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Order ID is required", common.ErrBadRequest, requestID)
		return
	}

	// Parse request
	var req UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request payload", err, map[string]interface{}{
			"request_id": requestID,
		})
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", common.ErrBadRequest, requestID)
		return
	}

	// Get existing order
	order, err := h.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		h.logger.Error("Failed to get order", err, map[string]interface{}{
			"request_id": requestID,
			"order_id":   orderID,
		})
		common.RespondWithError(w, http.StatusNotFound, "Order not found", common.ErrNotFound, requestID)
		return
	}

	// Update order status based on the request
	var eventType messaging.EventType
	switch domain.OrderStatus(req.Status) {
	case domain.OrderStatusConfirmed:
		if err := order.ConfirmOrder(); err != nil {
			common.RespondWithError(w, http.StatusBadRequest, err.Error(), common.ErrBadRequest, requestID)
			return
		}
		eventType = messaging.OrderConfirmed
	case domain.OrderStatusShipped:
		if err := order.ShipOrder(); err != nil {
			common.RespondWithError(w, http.StatusBadRequest, err.Error(), common.ErrBadRequest, requestID)
			return
		}
		eventType = messaging.OrderShipped
	case domain.OrderStatusDelivered:
		if err := order.DeliverOrder(); err != nil {
			common.RespondWithError(w, http.StatusBadRequest, err.Error(), common.ErrBadRequest, requestID)
			return
		}
	case domain.OrderStatusCancelled:
		if err := order.CancelOrder(); err != nil {
			common.RespondWithError(w, http.StatusBadRequest, err.Error(), common.ErrBadRequest, requestID)
			return
		}
	default:
		common.RespondWithError(w, http.StatusBadRequest, "Invalid order status", common.ErrBadRequest, requestID)
		return
	}

	// Save updated order
	if err := h.repo.UpdateOrder(ctx, order); err != nil {
		h.logger.Error("Failed to update order", err, map[string]interface{}{
			"request_id": requestID,
			"order_id":   orderID,
		})
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to update order", common.ErrInternalServerError, requestID)
		return
	}

	// Publish event if needed
	if eventType != "" {
		if err := h.publisher.Publish(eventType, order); err != nil {
			h.logger.Error("Failed to publish order event", err, map[string]interface{}{
				"request_id": requestID,
				"order_id":   orderID,
				"event_type": eventType,
			})
			// Continue even if the event publishing fails
		}
	}

	common.RespondWithJSON(w, http.StatusOK, order, requestID)
}

// CancelOrder handles DELETE /api/orders/{id}
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := chi.URLParam(r, "request_id")
	if requestID == "" {
		requestID = common.GenerateRequestID()
	}

	// Get order ID from URL
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Order ID is required", common.ErrBadRequest, requestID)
		return
	}

	// Get existing order
	order, err := h.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		h.logger.Error("Failed to get order", err, map[string]interface{}{
			"request_id": requestID,
			"order_id":   orderID,
		})
		common.RespondWithError(w, http.StatusNotFound, "Order not found", common.ErrNotFound, requestID)
		return
	}

	// Cancel the order
	if err := order.CancelOrder(); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, err.Error(), common.ErrBadRequest, requestID)
		return
	}

	// Save updated order
	if err := h.repo.UpdateOrder(ctx, order); err != nil {
		h.logger.Error("Failed to cancel order", err, map[string]interface{}{
			"request_id": requestID,
			"order_id":   orderID,
		})
		common.RespondWithError(w, http.StatusInternalServerError, "Failed to cancel order", common.ErrInternalServerError, requestID)
		return
	}

	common.RespondWithMessage(w, http.StatusOK, "Order cancelled successfully", requestID)
}

// checkInventory checks if a product is available in the inventory
func (h *OrderHandler) checkInventory(ctx context.Context, productID string, quantity int) (bool, error) {
	url := fmt.Sprintf("%s/api/inventory/%s?required=%d", h.inventoryURL, productID, quantity)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := h.inventoryClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

// publishOrderCreatedEvent publishes an event when an order is created
func (h *OrderHandler) publishOrderCreatedEvent(order *domain.Order) error {
	return h.publisher.Publish(messaging.OrderCreated, order)
}
