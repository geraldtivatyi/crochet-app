package domain

import "time"

type Order struct {
	ID            int       `json:"id"`
	ProductID     string    `json:"product_id"`
	Quantity      int       `json:"quantity"`
	CustomerEmail string    `json:"customer_email"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateOrderRequest struct {
	ProductID     string `json:"product_id"`
	Quantity      int    `json:"quantity"`
	CustomerEmail string `json:"customer_email"`
}
