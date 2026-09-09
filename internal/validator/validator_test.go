package validator_test

import (
	"testing"

	"Crochet/internal/domain"
	"Crochet/internal/validator"
)

func TestValidateProduct(t *testing.T) {
	tests := []struct {
		name        string
		product     domain.Product
		expectValid bool
		failedField string
	}{
		{
			name: "Valid Product",
			product: domain.Product{
				ID:           "prod_001",
				ProductName:  "Crochet Amigurumi Bear",
				ProductPrice: 150.00,
				StoreID:      "store_123",
			},
			expectValid: true,
		},
		{
			name: "Invalid - Empty Name",
			product: domain.Product{
				ID:           "prod_002",
				ProductName:  "   ",
				ProductPrice: 150.00,
			},
			expectValid: false,
			failedField: "product_name",
		},
		{
			name: "Invalid - Zero or Negative Price",
			product: domain.Product{
				ID:           "prod_003",
				ProductName:  "Crochet Hat",
				ProductPrice: -10.00,
			},
			expectValid: false,
			failedField: "product_price",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			validator.ValidateProduct(v, tt.product)

			if v.Valid() != tt.expectValid {
				t.Errorf("expected valid=%v, got %v", tt.expectValid, v.Valid())
			}

			if !tt.expectValid {
				if _, exists := v.Errors[tt.failedField]; !exists {
					t.Errorf("expected error field %q, but it was not found in %v", tt.failedField, v.Errors)
				}
			}
		})
	}
}
