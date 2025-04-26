package service

import (
	"context"
	"errors"
	
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/aggregate"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/entity"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/repository"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/shared/domain/events"
)

// OrderService provides domain logic for working with orders
type OrderService struct {
	orderRepository repository.OrderRepository
	eventBus        events.EventBus
}

// NewOrderService creates a new OrderService
func NewOrderService(orderRepo repository.OrderRepository, eventBus events.EventBus) *OrderService {
	return &OrderService{
		orderRepository: orderRepo,
		eventBus:        eventBus,
	}
}

// AddOrderItem adds an item to an order
func (s *OrderService) AddOrderItem(ctx context.Context, orderID string, item *entity.OrderItem) error {
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	
	if order == nil {
		return errors.New("order not found")
	}
	
	if err := order.AddItem(item); err != nil {
		return err
	}
	
	// Persist changes
	if err := s.orderRepository.Update(ctx, order); err != nil {
		return err
	}
	
	// Publish domain events
	for _, event := range order.Events() {
		s.eventBus.Publish(ctx, event)
	}
	
	// Clear events after publishing
	order.ClearEvents()
	
	return nil
}

// RemoveOrderItem removes an item from an order
func (s *OrderService) RemoveOrderItem(ctx context.Context, orderID, itemID string) error {
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	
	if order == nil {
		return errors.New("order not found")
	}
	
	if err := order.RemoveItem(itemID); err != nil {
		return err
	}
	
	// Persist changes
	if err := s.orderRepository.Update(ctx, order); err != nil {
		return err
	}
	
	// Publish domain events
	for _, event := range order.Events() {
		s.eventBus.Publish(ctx, event)
	}
	
	// Clear events after publishing
	order.ClearEvents()
	
	return nil
}

// CompleteOrderPayment marks an order as paid
func (s *OrderService) CompleteOrderPayment(ctx context.Context, orderID string) error {
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	
	if order == nil {
		return errors.New("order not found")
	}
	
	if err := order.MarkAsPaid(); err != nil {
		return err
	}
	
	// Persist changes
	if err := s.orderRepository.Update(ctx, order); err != nil {
		return err
	}
	
	// Publish domain events
	for _, event := range order.Events() {
		s.eventBus.Publish(ctx, event)
	}
	
	// Clear events after publishing
	order.ClearEvents()
	
	return nil
}

// ShipOrder marks an order as shipped
func (s *OrderService) ShipOrder(ctx context.Context, orderID string) error {
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	
	if order == nil {
		return errors.New("order not found")
	}
	
	if err := order.MarkAsShipped(); err != nil {
		return err
	}
	
	// Persist changes
	if err := s.orderRepository.Update(ctx, order); err != nil {
		return err
	}
	
	// Publish domain events
	for _, event := range order.Events() {
		s.eventBus.Publish(ctx, event)
	}
	
	// Clear events after publishing
	order.ClearEvents()
	
	return nil
}

// DeliverOrder marks an order as delivered
func (s *OrderService) DeliverOrder(ctx context.Context, orderID string) error {
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	
	if order == nil {
		return errors.New("order not found")
	}
	
	if err := order.MarkAsDelivered(); err != nil {
		return err
	}
	
	// Persist changes
	if err := s.orderRepository.Update(ctx, order); err != nil {
		return err
	}
	
	// Publish domain events
	for _, event := range order.Events() {
		s.eventBus.Publish(ctx, event)
	}
	
	// Clear events after publishing
	order.ClearEvents()
	
	return nil
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(ctx context.Context, orderID string) error {
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	
	if order == nil {
		return errors.New("order not found")
	}
	
	if err := order.Cancel(); err != nil {
		return err
	}
	
	// Persist changes
	if err := s.orderRepository.Update(ctx, order); err != nil {
		return err
	}
	
	// Publish domain events
	for _, event := range order.Events() {
		s.eventBus.Publish(ctx, event)
	}
	
	// Clear events after publishing
	order.ClearEvents()
	
	return nil
}
