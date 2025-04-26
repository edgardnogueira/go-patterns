package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/edgardnogueira/go-patterns/project_structures/microservices/inventory-service/internal/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// InventoryRepository defines the interface for inventory data access
type InventoryRepository interface {
	GetProductByID(ctx context.Context, id string) (*domain.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error)
	ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error)
	CreateProduct(ctx context.Context, product *domain.Product) error
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, id string) error
	CreateReservation(ctx context.Context, reservation *domain.InventoryReservation) error
	GetReservationsByOrderID(ctx context.Context, orderID string) ([]*domain.InventoryReservation, error)
	UpdateReservationStatus(ctx context.Context, id, status string) error
	Close() error
}

// PostgresInventoryRepository implements InventoryRepository using PostgreSQL
type PostgresInventoryRepository struct {
	db *sqlx.DB
}

// NewPostgresInventoryRepository creates a new repository
func NewPostgresInventoryRepository(connStr string) (*PostgresInventoryRepository, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Setup database
	if err := setupDatabase(db); err != nil {
		db.Close()
		return nil, err
	}

	return &PostgresInventoryRepository{db: db}, nil
}

// setupDatabase creates the necessary tables if they don't exist
func setupDatabase(db *sqlx.DB) error {
	// Create products table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(64) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			sku VARCHAR(64) UNIQUE NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			quantity INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create inventory_reservations table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS inventory_reservations (
			id VARCHAR(64) PRIMARY KEY,
			product_id VARCHAR(64) NOT NULL,
			order_id VARCHAR(64) NOT NULL,
			quantity INTEGER NOT NULL,
			reserved_at TIMESTAMP NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			status VARCHAR(20) NOT NULL,
			FOREIGN KEY (product_id) REFERENCES products(id)
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// GetProductByID retrieves a product by its ID
func (r *PostgresInventoryRepository) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	var product domain.Product

	query := `
		SELECT id, name, description, sku, price, quantity, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		return nil, err
	}

	return &product, nil
}

// GetProductBySKU retrieves a product by its SKU
func (r *PostgresInventoryRepository) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	var product domain.Product

	query := `
		SELECT id, name, description, sku, price, quantity, created_at, updated_at
		FROM products
		WHERE sku = $1
	`
	err := r.db.GetContext(ctx, &product, query, sku)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		return nil, err
	}

	return &product, nil
}

// ListProducts retrieves a list of products with pagination
func (r *PostgresInventoryRepository) ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	var products []*domain.Product

	query := `
		SELECT id, name, description, sku, price, quantity, created_at, updated_at
		FROM products
		ORDER BY name
		LIMIT $1 OFFSET $2
	`
	err := r.db.SelectContext(ctx, &products, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// CreateProduct creates a new product
func (r *PostgresInventoryRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (id, name, description, sku, price, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Description,
		product.SKU,
		product.Price,
		product.Quantity,
		product.CreatedAt,
		product.UpdatedAt,
	)
	return err
}

// UpdateProduct updates an existing product
func (r *PostgresInventoryRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products 
		SET name = $1, description = $2, sku = $3, price = $4, quantity = $5, updated_at = $6
		WHERE id = $7
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.SKU,
		product.Price,
		product.Quantity,
		product.UpdatedAt,
		product.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

// DeleteProduct deletes a product
func (r *PostgresInventoryRepository) DeleteProduct(ctx context.Context, id string) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

// CreateReservation creates a new inventory reservation
func (r *PostgresInventoryRepository) CreateReservation(ctx context.Context, reservation *domain.InventoryReservation) error {
	query := `
		INSERT INTO inventory_reservations (id, product_id, order_id, quantity, reserved_at, expires_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		reservation.ID,
		reservation.ProductID,
		reservation.OrderID,
		reservation.Quantity,
		reservation.ReservedAt,
		reservation.ExpiresAt,
		reservation.Status,
	)
	return err
}

// GetReservationsByOrderID retrieves all reservations for an order
func (r *PostgresInventoryRepository) GetReservationsByOrderID(ctx context.Context, orderID string) ([]*domain.InventoryReservation, error) {
	var reservations []*domain.InventoryReservation

	query := `
		SELECT id, product_id, order_id, quantity, reserved_at, expires_at, status
		FROM inventory_reservations
		WHERE order_id = $1
	`
	err := r.db.SelectContext(ctx, &reservations, query, orderID)
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

// UpdateReservationStatus updates the status of a reservation
func (r *PostgresInventoryRepository) UpdateReservationStatus(ctx context.Context, id, status string) error {
	query := `
		UPDATE inventory_reservations 
		SET status = $1
		WHERE id = $2
	`
	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("reservation not found")
	}

	return nil
}

// Close closes the database connection
func (r *PostgresInventoryRepository) Close() error {
	return r.db.Close()
}
