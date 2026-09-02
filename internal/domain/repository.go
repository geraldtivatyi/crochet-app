package domain

import "context"

type ProductRepository interface {
	GetProducts(ctx context.Context, storeID string) ([]Product, error)
	AddProduct(ctx context.Context, p Product) error
	DeleteProduct(ctx context.Context, id string) error
	UpdateProduct(ctx context.Context, id string, p Product) error
	Close() error
}
