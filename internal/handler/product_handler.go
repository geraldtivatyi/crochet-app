package handler

import (
	"Crochet/internal/auth"
	"Crochet/internal/queue"
	"Crochet/internal/response"
	"Crochet/internal/validator"
	"encoding/json"
	"fmt"
	"net/http"

	"Crochet/internal/domain"
)

func HandleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func HandleGetProducts(repo domain.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storeID := r.URL.Query().Get("storeID")

		// If query param isn't set, try to get storeID from authenticated context
		if storeID == "" {
			if user, err := auth.UserFromContext(r.Context()); err == nil {
				storeID = user.StoreID
			}
		}

		products, err := repo.GetProducts(r.Context(), storeID)
		if err != nil {
			http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, products)
	}
}

func HandleAddProduct(repo domain.ProductRepository, taskQ *queue.TaskQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.UserFromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}

		var p domain.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, fmt.Sprintf("Invalid payload: %v", err), http.StatusBadRequest)
			return
		}

		p.StoreID = user.StoreID

		// 1. Validate payload fields
		v := validator.New()
		if validator.ValidateProduct(v, p); !v.Valid() {
			response.ValidationError(w, v.Errors)
			return
		}

		// 2. Save valid product to store
		if err := repo.AddProduct(r.Context(), p); err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to save product to database")
			return
		}

		taskQ.Enqueue(queue.NotificationTask{
			ProductID: p.ID,
			Action:    "CREATED",
		})

		response.JSON(w, http.StatusCreated, p)
	}
}

func HandleDeleteProduct(repo domain.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "Missing product ID", http.StatusBadRequest)
			return
		}

		if err := repo.DeleteProduct(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleUpdateProduct(repo domain.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.UserFromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}

		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "Missing product ID", http.StatusBadRequest)
			return
		}

		var p domain.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, fmt.Sprintf("Invalid payload: %v", err), http.StatusBadRequest)
			return
		}

		// Ensure the struct ID matches the URL path ID
		p.ID = id
		p.StoreID = user.StoreID

		if err := repo.UpdateProduct(r.Context(), id, p); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}
