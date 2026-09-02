package handler_test

import (
	"Crochet/internal/domain"
	"context"
)

// MockStore implements domain.ProductRepository for testing
type MockStore struct {
	products         []domain.Product
	addErr           error
	getErr           error
	addProductCalled bool
}

func (m *MockStore) GetProducts(ctx context.Context) ([]domain.Product, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.products, nil
}

func (m *MockStore) AddProduct(ctx context.Context, p domain.Product) error {
	m.addProductCalled = true
	if m.addErr != nil {
		return m.addErr
	}
	m.products = append(m.products, p)
	return nil
}

func (m *MockStore) Close() error {
	return nil
}
