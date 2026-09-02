package store

import (
	"context"
	"database/sql"
	"fmt"

	"Crochet/internal/domain"

	_ "github.com/glebarez/go-sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	// Create table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY,
		product_name TEXT NOT NULL,
		product_price REAL NOT NULL,
		store_id TEXT NOT NULL
	);`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create products table: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) GetProducts(ctx context.Context, storeID string) ([]domain.Product, error) {
	queryStr := "SELECT id, product_name, product_price, store_id FROM products"
	var args []any

	if storeID != "" {
		queryStr += ` WHERE store_id = ?`
		args = append(args, storeID)
	}

	rows, err := s.db.QueryContext(ctx, queryStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.ProductName, &p.ProductPrice, &p.StoreID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (s *SQLiteStore) AddProduct(ctx context.Context, p domain.Product) error {
	query := `INSERT INTO products (id, product_name, product_price, store_id) VALUES (?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, p.ID, p.ProductName, p.ProductPrice, p.StoreID)
	return err
}

func (s *SQLiteStore) DeleteProduct(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %s not found", id)
	}

	return nil
}

func (s *SQLiteStore) UpdateProduct(ctx context.Context, id string, p domain.Product) error {
	query := `
	UPDATE products 
	SET product_name = ?, product_price = ?, store_id = ? 
	WHERE id = ?`

	res, err := s.db.ExecContext(ctx, query, p.ProductName, p.ProductPrice, p.StoreID, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %s not found", id)
	}

	return nil
}

func (s *SQLiteStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
