package mapper

import (
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/application/dto"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/aggregate"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/entity"
	"github.com/edgardnogueira/go-patterns/project_structures/ddd/internal/order/domain/valueobject"
	"github.com/google/uuid"
)

// OrderMapper provides mapping functions between domain objects and DTOs
type OrderMapper struct{}

// NewOrderMapper creates a new OrderMapper
func NewOrderMapper() *OrderMapper {
	return &OrderMapper{}
}

// ToOrderDTO converts an order aggregate to an order DTO
func (m *OrderMapper) ToOrderDTO(order *aggregate.Order) (*dto.OrderDTO, error) {
	if order == nil {
		return nil, nil
	}
	
	// Convert order items
	var itemDTOs []dto.OrderItemDTO
	for _, item := range order.Items() {
		itemDTO := dto.OrderItemDTO{
			ID:          item.ID(),
			ProductID:   item.ProductID(),
			Quantity:    item.Quantity(),
			UnitPrice:   float64(item.UnitPrice().Amount()) / 100, // Convert cents to dollars
			Currency:    item.UnitPrice().Currency(),
			Description: item.Description(),
		}
		itemDTOs = append(itemDTOs, itemDTO)
	}
	
	// Calculate total price
	total, err := order.CalculateTotal()
	if err != nil {
		return nil, err
	}
	
	// Create DTO
	orderDTO := &dto.OrderDTO{
		ID:         order.ID(),
		CustomerID: order.CustomerID(),
		Items:      itemDTOs,
		Status:     dto.OrderStatus(order.Status()),
		TotalPrice: float64(total.Amount()) / 100, // Convert cents to dollars
		Currency:   total.Currency(),
		CreatedAt:  order.CreatedAt(),
		UpdatedAt:  order.UpdatedAt(),
	}
	
	return orderDTO, nil
}

// ToOrderEntity converts an AddOrderItemRequest DTO to an OrderItem entity
func (m *OrderMapper) ToOrderItemEntity(request *dto.AddOrderItemRequest) (*entity.OrderItem, error) {
	// Convert price from dollars to cents
	amountInCents := int64(request.UnitPrice * 100)
	
	// Create money value object
	money, err := valueobject.NewMoney(amountInCents, request.Currency)
	if err != nil {
		return nil, err
	}
	
	// Create entity
	return entity.NewOrderItem(
		uuid.New().String(),
		request.ProductID,
		request.Quantity,
		money,
		request.Description,
	)
}

// ToOrderAggregate creates a new Order aggregate from a CreateOrderRequest
func (m *OrderMapper) ToOrderAggregate(request *dto.CreateOrderRequest) (*aggregate.Order, error) {
	return aggregate.NewOrder(
		uuid.New().String(),
		request.CustomerID,
	)
}
