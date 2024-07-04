package postgres

import (
	"database/sql"
	"errors"
	"product-service/models"

	"github.com/google/uuid"
)

var ErrProductNotFound = errors.New("user not found")

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) AddProduct(product *models.Product) (*models.Product, error) {
	id := uuid.New()
	query := `INSERT INTO products (id, name, description, price, stock, created_at::text, updated_at::text) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, name, description, price, stock, created_at, updated_at`
	err := r.db.QueryRow(query, id, product.Name, product.Description, product.Price, product.Stock).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	return product, err
}

func (r *PostgresRepository) GetProduct(id string) (*models.Product, error) {
	query := `SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var product models.Product
	if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *PostgresRepository) UpdateProduct(product *models.Product) (*models.Product, error) {
	query := `UPDATE products SET name = $2, description = $3, price = $4, stock = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING id, name, description, price, stock, created_at, updated_at`
	err := r.db.QueryRow(query, product.ID, product.Name, product.Description, product.Price, product.Stock).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	return product, err
}

func (r *PostgresRepository) DeleteProduct(id string) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (r *PostgresRepository) ListProducts(limit, page int) ([]*models.Product, error) {
	if page == 1 {
		page = 0
	}
	query := `SELECT id, name, description, price, stock, created_at, updated_at FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}
