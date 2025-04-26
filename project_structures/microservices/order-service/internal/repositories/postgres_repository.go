package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/edgardnogueira/go-patterns/project_structures/microservices/order-service/internal/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrderByID(ctx context.Context, id string) (*domain.Order, error)
	ListOrders(ctx context.Context, limit, offset int) ([]*domain.Order, error)
	UpdateOrder(ctx context.Context, order *domain.Order) error
	DeleteOrder(ctx context.Context, id string) error
	AddOrderItem(ctx context.Context, item *domain.OrderItem) error
	GetOrderItems(ctx context.Context, orderID string) ([]domain.OrderItem, error)
	UpdateOrderItem(ctx context.Context, item *domain.OrderItem) error
	DeleteOrderItem(ctx context.Context, id string) error
	Close() error
}

// PostgresOrderRepository implements OrderRepository using PostgreSQL
type PostgresOrderRepository struct {
	db *sqlx.DB
}

// NewPostgresOrderRepository creates a new repository
func NewPostgresOrderRepository(connStr string) (*PostgresOrderRepository, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Setup database
	if err := setupDatabase(db); err != nil {
		db.Close()
		return nil, err
	}

	return &PostgresOrderRepository{db: db}, nil
}

// setupDatabase creates the necessary tables if they don't exist
func setupDatabase(db *sqlx.DB) error {
	// Create orders table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id VARCHAR(64) PRIMARY KEY,
			customer_id VARCHAR(64) NOT NULL,
			status VARCHAR(20) NOT NULL,
			total_amount DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			shipped_at TIMESTAMP,
			delivered_at TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create order_items table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
			id VARCHAR(64) PRIMARY KEY,
			order_id VARCHAR(64) NOT NULL,
			product_id VARCHAR(64) NOT NULL,
			quantity INTEGER NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// CreateOrder saves a new order to the database
func (r *PostgresOrderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert order
	query := `
		INSERT INTO orders (id, customer_id, status, total_amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.ExecContext(
		ctx,
		query,
		order.ID,
		order.CustomerID,
		order.Status,
		order.TotalAmount,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert order items
	for _, item := range order.Items {
		query := `
			INSERT INTO order_items (id, order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err = tx.ExecContext(
			ctx,
			query,
			item.ID,
			order.ID,
			item.ProductID,
			item.Quantity,
			item.Price,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetOrderByID retrieves an order by its ID
func (r *PostgresOrderRepository) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order

	query := `
		SELECT id, customer_id, status, total_amount, created_at, updated_at, shipped_at, delivered_at
		FROM orders
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &order, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found: %w", err)
		}
		return nil, err
	}

	// Get order items
	items, err := r.GetOrderItems(ctx, id)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return &order, nil
}

// ListOrders retrieves a list of orders with pagination
func (r *PostgresOrderRepository) ListOrders(ctx context.Context, limit, offset int) ([]*domain.Order, error) {
	var orders []*domain.Order

	query := `
		SELECT id, customer_id, status, total_amount, created_at, updated_at, shipped_at, delivered_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	err := r.db.SelectContext(ctx, &orders, query, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get items for each order
	for i, order := range orders {
		items, err := r.GetOrderItems(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}

	return orders, nil
}

// UpdateOrder updates an existing order
func (r *PostgresOrderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	query := `
		UPDATE orders 
		SET status = $1, total_amount = $2, updated_at = $3, shipped_at = $4, delivered_at = $5
		WHERE id = $6
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		order.Status,
		order.TotalAmount,
		order.UpdatedAt,
		order.ShippedAt,
		order.DeliveredAt,
		order.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil
}

// DeleteOrder deletes an order and its items
func (r *PostgresOrderRepository) DeleteOrder(ctx context.Context, id string) error {
	// The database has ON DELETE CASCADE for order_items
	query := "DELETE FROM orders WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil
}

// AddOrderItem adds a new item to an order
func (r *PostgresOrderRepository) AddOrderItem(ctx context.Context, item *domain.OrderItem) error {
	query := `
		INSERT INTO order_items (id, order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		item.ID,
		item.OrderID,
		item.ProductID,
		item.Quantity,
		item.Price,
	)
	return err
}

// GetOrderItems retrieves all items for an order
func (r *PostgresOrderRepository) GetOrderItems(ctx context.Context, orderID string) ([]domain.OrderItem, error) {
	var items []domain.OrderItem

	query := `
		SELECT id, order_id, product_id, quantity, price
		FROM order_items
		WHERE order_id = $1
	`
	err := r.db.SelectContext(ctx, &items, query, orderID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// UpdateOrderItem updates an existing order item
func (r *PostgresOrderRepository) UpdateOrderItem(ctx context.Context, item *domain.OrderItem) error {
	query := `
		UPDATE order_items 
		SET quantity = $1, price = $2
		WHERE id = $3
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		item.Quantity,
		item.Price,
		item.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("order item not found")
	}

	return nil
}

// DeleteOrderItem deletes an order item
func (r *PostgresOrderRepository) DeleteOrderItem(ctx context.Context, id string) error {
	query := "DELETE FROM order_items WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("order item not found")
	}

	return nil
}

// Close closes the database connection
func (r *PostgresOrderRepository) Close() error {
	return r.db.Close()
}
