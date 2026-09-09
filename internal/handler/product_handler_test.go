package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"Crochet/internal/auth"
	"Crochet/internal/domain"
	"Crochet/internal/handler"
	"Crochet/internal/queue"
)

// MockRepository implements domain.ProductRepository for testing
type MockRepository struct {
	products  []domain.Product
	addCalled bool
}

func (m *MockRepository) GetProducts(ctx context.Context, storeID string) ([]domain.Product, error) {
	return m.products, nil
}

func (m *MockRepository) AddProduct(ctx context.Context, p domain.Product) error {
	m.addCalled = true
	m.products = append(m.products, p)
	return nil
}

func (m *MockRepository) DeleteProduct(ctx context.Context, id string) error { return nil }
func (m *MockRepository) UpdateProduct(ctx context.Context, id string, p domain.Product) error {
	return nil
}
func (m *MockRepository) Close() error { return nil }

func TestHandleAddProduct_ValidationFailure(t *testing.T) {
	mockRepo := &MockRepository{}
	taskQ := queue.NewTaskQueue(10, 1)

	// Payload with missing name and negative price
	invalidPayload := `{"id":"p101","product_name":"","product_price":-5.00}`

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(invalidPayload))
	req.Header.Set("Content-Type", "application/json")

	// Inject authenticated user into context manually for testing
	ctx := auth.WithUser(req.Context(), auth.UserClaims{
		UserID:  "user_1",
		StoreID: "store_cape_town",
	})
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	httpHandler := handler.HandleAddProduct(mockRepo, taskQ)
	httpHandler.ServeHTTP(rec, req)

	// 1. Assert HTTP status is 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	// 2. Assert database was NOT called
	if mockRepo.addCalled {
		t.Errorf("expected AddProduct not to be called on database for invalid payload")
	}
}
