package http

import (
	"encoding/json"
	"net/http"
	
	"github.com/go-chi/chi/v5"
	
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/model"
	"github.com/edgardnogueira/go-patterns/project_structures/hexagonal/internal/domain/service"
)

// OrderHandler is the HTTP adapter for order operations
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler creates a new HTTP handler for order operations
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// RegisterRoutes registers the HTTP routes for order operations
func (h *OrderHandler) RegisterRoutes(r chi.Router) {
	r.Post("/orders", h.CreateOrder)
	r.Get("/orders/{orderID}", h.GetOrder)
	r.Put("/orders/{orderID}/status", h.UpdateOrderStatus)
	r.Get("/customers/{customerID}/orders", h.ListOrders)
	r.Post("/orders/{orderID}/items", h.AddItemToOrder)
	r.Delete("/orders/{orderID}/items/{productID}", h.RemoveItemFromOrder)
}

// CreateOrderRequest represents the payload for creating a new order
type CreateOrderRequest struct {
	CustomerID      string         `json:"customer_id"`
	Items           []model.OrderItem `json:"items"`
	ShippingAddress string         `json:"shipping_address"`
}

// UpdateOrderStatusRequest represents the payload for updating an order's status
type UpdateOrderStatusRequest struct {
	Status model.OrderStatus `json:"status"`
}

// AddItemRequest represents the payload for adding an item to an order
type AddItemRequest struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// CreateOrder handles the HTTP request to create a new order
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	
	order, err := h.orderService.CreateOrder(r.Context(), req.CustomerID, req.Items, req.ShippingAddress)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusCreated, order)
}

// GetOrder handles the HTTP request to retrieve an order by ID
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")
	if orderID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing order ID")
		return
	}
	
	order, err := h.orderService.GetOrder(r.Context(), orderID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}
	
	respondWithJSON(w, http.StatusOK, order)
}

// UpdateOrderStatus handles the HTTP request to update an order's status
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")
	if orderID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing order ID")
		return
	}
	
	var req UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	
	if err := h.orderService.UpdateOrderStatus(r.Context(), orderID, req.Status); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// ListOrders handles the HTTP request to list orders for a customer
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	customerID := chi.URLParam(r, "customerID")
	if customerID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing customer ID")
		return
	}
	
	orders, err := h.orderService.ListOrders(r.Context(), customerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, orders)
}

// AddItemToOrder handles the HTTP request to add an item to an order
func (h *OrderHandler) AddItemToOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")
	if orderID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing order ID")
		return
	}
	
	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	
	item := model.OrderItem{
		ProductID: req.ProductID,
		Name:      req.Name,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
	}
	
	if err := h.orderService.AddItemToOrder(r.Context(), orderID, item); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "item added"})
}

// RemoveItemFromOrder handles the HTTP request to remove an item from an order
func (h *OrderHandler) RemoveItemFromOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")
	productID := chi.URLParam(r, "productID")
	
	if orderID == "" || productID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing order ID or product ID")
		return
	}
	
	if err := h.orderService.RemoveItemFromOrder(r.Context(), orderID, productID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "item removed"})
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
