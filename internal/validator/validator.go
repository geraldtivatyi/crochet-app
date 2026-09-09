package validator

import (
	"strings"

	"Crochet/internal/domain"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if no validation rules failed
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError records a validation failure for a specific field if it hasn't been set yet
func (v *Validator) AddError(field, message string) {
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}

// Check adds an error if the condition evaluates to false
func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

// ValidateProduct checks all business constraints for a Product struct
func ValidateProduct(v *Validator, p domain.Product) {
	v.Check(strings.TrimSpace(p.ID) != "", "id", "Product ID cannot be empty")
	v.Check(strings.TrimSpace(p.ProductName) != "", "product_name", "Product name cannot be empty")
	v.Check(len(p.ProductName) <= 100, "product_name", "Product name must not exceed 100 characters")
	v.Check(p.ProductPrice > 0, "product_price", "Product price must be greater than zero")
}

// ValidateOrderRequest checks business rules for incoming orders
func ValidateOrderRequest(v *Validator, req domain.CreateOrderRequest) {
	v.Check(strings.TrimSpace(req.ProductID) != "", "product_id", "Product ID is required")
	v.Check(req.Quantity > 0, "quantity", "Quantity must be at least 1")
	v.Check(strings.TrimSpace(req.CustomerEmail) != "", "customer_email", "Customer email is required")
	v.Check(strings.Contains(req.CustomerEmail, "@"), "customer_email", "Invalid email address format")
}
