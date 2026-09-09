package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Crochet/internal/domain"
	"Crochet/internal/queue"
	"Crochet/internal/response"
	"Crochet/internal/validator"
)

func HandleCreateOrder(repo domain.ProductRepository, taskQ *queue.TaskQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON payload syntax")
			return
		}

		// 1. Validate payload fields
		v := validator.New()
		if validator.ValidateOrderRequest(v, req); !v.Valid() {
			response.ValidationError(w, v.Errors)
			return
		}

		// 2. Execute SQL Transaction (Deduct Stock + Create Order Record)
		err := repo.CreateOrderTx(r.Context(), req.ProductID, req.Quantity, req.CustomerEmail)
		if err != nil {
			// Check if failure was due to insufficient stock or database error
			response.Error(w, http.StatusBadRequest, fmt.Sprintf("Order processing failed: %v", err))
			return
		}

		// 3. Enqueue background notification to alert artisan of new sale
		taskQ.Enqueue(queue.NotificationTask{
			ProductID: req.ProductID,
			Action:    fmt.Sprintf("ORDER_PLACED_QTY_%d", req.Quantity),
		})

		// 4. Return success response
		response.JSON(w, http.StatusCreated, map[string]string{
			"message": "Order created successfully",
		})
	}
}
