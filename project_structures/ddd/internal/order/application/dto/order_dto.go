package dto

import (
	"time"
)

// OrderStatus represents the status of an order in DTO form
type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "CREATED"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// OrderItemDTO is a data transfer object for order items
type OrderItemDTO struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
}

// OrderDTO is a data transfer object for orders
type OrderDTO struct {
	ID         string         `json:"id"`
	CustomerID string         `json:"customer_id"`
	Items      []OrderItemDTO `json:"items"`
	Status     OrderStatus    `json:"status"`
	TotalPrice float64        `json:"total_price"`
	Currency   string         `json:"currency"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// CreateOrderRequest represents a request to create a new order
type CreateOrderRequest struct {
	CustomerID string `json:"customer_id"`
}

// AddOrderItemRequest represents a request to add an item to an order
type AddOrderItemRequest struct {
	ProductID   string  `json:"product_id"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
}

// PayOrderRequest represents a request to pay for an order
type PayOrderRequest struct {
	PaymentMethod string `json:"payment_method"`
	PaymentAmount float64 `json:"payment_amount"`
	Currency      string `json:"currency"`
}
