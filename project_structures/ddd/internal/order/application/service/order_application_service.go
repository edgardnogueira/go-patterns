package service

import (
	"context"
	"errors"
	
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/application/dto"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/application/mapper"
	domainService "github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/service"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/repository"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/shared/domain/events"
)

// OrderApplicationService handles use cases related to orders
type OrderApplicationService struct {
	orderRepository repository.OrderRepository
	orderService    *domainService.OrderService
	orderMapper     *mapper.OrderMapper
	eventBus        events.EventBus
}

// NewOrderApplicationService creates a new OrderApplicationService
func NewOrderApplicationService(
	orderRepo repository.OrderRepository,
	eventBus events.EventBus,
) *OrderApplicationService {
	orderService := domainService.NewOrderService(orderRepo, eventBus)
	orderMapper := mapper.NewOrderMapper()
	
	return &OrderApplicationService{
		orderRepository: orderRepo,
		orderService:    orderService,
		orderMapper:     orderMapper,
		eventBus:        eventBus,
	}
}

// CreateOrder creates a new order
func (s *OrderApplicationService) CreateOrder(ctx context.Context, request *dto.CreateOrderRequest) (*dto.OrderDTO, error) {
	// Validate request
	if request.CustomerID == "" {
		return nil, errors.New("customer ID is required")
	}
	
	// Create order aggregate
	order, err := s.orderMapper.ToOrderAggregate(request)
	if err != nil {
		return nil, err
	}
	
	// Save order
	if err := s.orderRepository.Save(ctx, order); err != nil {
		return nil, err
	}
	
	// Publish domain events
	for _, event := range order.Events() {
		s.eventBus.Publish(ctx, event)
	}
	
	// Clear events
	order.ClearEvents()
	
	// Map to DTO and return
	return s.orderMapper.ToOrderDTO(order)
}

// GetOrder retrieves an order by ID
func (s *OrderApplicationService) GetOrder(ctx context.Context, orderID string) (*dto.OrderDTO, error) {
	// Retrieve order from repository
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	
	if order == nil {
		return nil, errors.New("order not found")
	}
	
	// Map to DTO and return
	return s.orderMapper.ToOrderDTO(order)
}

// GetOrdersByCustomer retrieves all orders for a customer
func (s *OrderApplicationService) GetOrdersByCustomer(ctx context.Context, customerID string) ([]*dto.OrderDTO, error) {
	// Retrieve orders from repository
	orders, err := s.orderRepository.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	
	// Map to DTOs
	orderDTOs := make([]*dto.OrderDTO, 0, len(orders))
	for _, order := range orders {
		orderDTO, err := s.orderMapper.ToOrderDTO(order)
		if err != nil {
			return nil, err
		}
		orderDTOs = append(orderDTOs, orderDTO)
	}
	
	return orderDTOs, nil
}

// AddOrderItem adds an item to an order
func (s *OrderApplicationService) AddOrderItem(ctx context.Context, orderID string, request *dto.AddOrderItemRequest) (*dto.OrderDTO, error) {
	// Convert DTO to domain entity
	item, err := s.orderMapper.ToOrderItemEntity(request)
	if err != nil {
		return nil, err
	}
	
	// Add item to order (using domain service)
	if err := s.orderService.AddOrderItem(ctx, orderID, item); err != nil {
		return nil, err
	}
	
	// Get updated order
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	
	// Map to DTO and return
	return s.orderMapper.ToOrderDTO(order)
}

// PayOrder processes payment for an order
func (s *OrderApplicationService) PayOrder(ctx context.Context, orderID string, request *dto.PayOrderRequest) (*dto.OrderDTO, error) {
	// In a real implementation, we would process the payment here
	// For this example, we'll just mark the order as paid
	
	// Mark order as paid (using domain service)
	if err := s.orderService.CompleteOrderPayment(ctx, orderID); err != nil {
		return nil, err
	}
	
	// Get updated order
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	
	// Map to DTO and return
	return s.orderMapper.ToOrderDTO(order)
}

// ShipOrder marks an order as shipped
func (s *OrderApplicationService) ShipOrder(ctx context.Context, orderID string) (*dto.OrderDTO, error) {
	// Mark order as shipped (using domain service)
	if err := s.orderService.ShipOrder(ctx, orderID); err != nil {
		return nil, err
	}
	
	// Get updated order
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	
	// Map to DTO and return
	return s.orderMapper.ToOrderDTO(order)
}

// DeliverOrder marks an order as delivered
func (s *OrderApplicationService) DeliverOrder(ctx context.Context, orderID string) (*dto.OrderDTO, error) {
	// Mark order as delivered (using domain service)
	if err := s.orderService.DeliverOrder(ctx, orderID); err != nil {
		return nil, err
	}
	
	// Get updated order
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	
	// Map to DTO and return
	return s.orderMapper.ToOrderDTO(order)
}

// CancelOrder cancels an order
func (s *OrderApplicationService) CancelOrder(ctx context.Context, orderID string) (*dto.OrderDTO, error) {
	// Cancel order (using domain service)
	if err := s.orderService.CancelOrder(ctx, orderID); err != nil {
		return nil, err
	}
	
	// Get updated order
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	
	// Map to DTO and return
	return s.orderMapper.ToOrderDTO(order)
}
