package builder

import (
	"database/sql"
	"fmt"

	"github.com/edgardnogueira/go-patterns/project_structures/event_driven/internal/events/types"
	_ "github.com/lib/pq"
)

// ProductProjectionBuilder builds and maintains product projections
type ProductProjectionBuilder struct {
	db *sql.DB
}

// NewProductProjectionBuilder creates a new ProductProjectionBuilder
func NewProductProjectionBuilder(connectionString string) *ProductProjectionBuilder {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping database: %v", err))
	}

	// Ensure tables exist
	if err = createProjectionTables(db); err != nil {
		panic(fmt.Sprintf("Failed to create projection tables: %v", err))
	}

	return &ProductProjectionBuilder{db: db}
}

// createProjectionTables creates the necessary tables for projections
func createProjectionTables(db *sql.DB) error {
	// Create products table for read model
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products_view (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10, 2) NOT NULL,
			stock INT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)

	return err
}

// CreateProduct creates a new product in the projection
func (b *ProductProjectionBuilder) CreateProduct(event types.ProductCreatedEvent) error {
	// Insert into products_view table
	_, err := b.db.Exec(
		`INSERT INTO products_view (id, name, description, price, stock, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $6)`,
		event.ID,
		event.Name,
		event.Description,
		event.Price,
		event.Stock,
		event.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to insert product projection: %w", err)
	}

	return nil
}

// UpdateProductStock updates the stock level of a product in the projection
func (b *ProductProjectionBuilder) UpdateProductStock(productID string, newStock int) error {
	// Update the product in the projection
	res, err := b.db.Exec(
		`UPDATE products_view SET stock = $1, updated_at = NOW() WHERE id = $2`,
		newStock,
		productID,
	)

	if err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %s not found", productID)
	}

	return nil
}

// GetProduct retrieves a product from the projection
func (b *ProductProjectionBuilder) GetProduct(productID string) (map[string]interface{}, error) {
	var id, name, description string
	var price float64
	var stock int
	var createdAt, updatedAt sql.NullTime

	err := b.db.QueryRow(
		"SELECT id, name, description, price, stock, created_at, updated_at FROM products_view WHERE id = $1",
		productID,
	).Scan(&id, &name, &description, &price, &stock, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product with ID %s not found", productID)
	}

	if err != nil {
		return nil, fmt.Errorf("error querying product: %w", err)
	}

	product := map[string]interface{}{
		"id":          id,
		"name":        name,
		"description": description,
		"price":       price,
		"stock":       stock,
	}

	if createdAt.Valid {
		product["created_at"] = createdAt.Time
	}

	if updatedAt.Valid {
		product["updated_at"] = updatedAt.Time
	}

	return product, nil
}

// GetProductsWithLowStock retrieves all products with stock below a threshold
func (b *ProductProjectionBuilder) GetProductsWithLowStock(threshold int) ([]map[string]interface{}, error) {
	rows, err := b.db.Query(
		"SELECT id, name, description, price, stock FROM products_view WHERE stock <= $1",
		threshold,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying products with low stock: %w", err)
	}
	defer rows.Close()

	var products []map[string]interface{}

	for rows.Next() {
		var id, name, description string
		var price float64
		var stock int

		if err := rows.Scan(&id, &name, &description, &price, &stock); err != nil {
			return nil, fmt.Errorf("error scanning product row: %w", err)
		}

		product := map[string]interface{}{
			"id":          id,
			"name":        name,
			"description": description,
			"price":       price,
			"stock":       stock,
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	return products, nil
}

// Close closes the database connection
func (b *ProductProjectionBuilder) Close() error {
	return b.db.Close()
}
