package handler

import (
	"Crochet/internal/queue"
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

		products, err := repo.GetProducts(r.Context(), storeID)
		if err != nil {
			http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func HandleAddProduct(repo domain.ProductRepository, taskQ *queue.TaskQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p domain.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, fmt.Sprintf("Invalid payload: %v", err), http.StatusBadRequest)
			return
		}

		if err := repo.AddProduct(r.Context(), p); err != nil {
			http.Error(w, "Failed to save product", http.StatusInternalServerError)
			return
		}

		taskQ.Enqueue(queue.NotificationTask{
			ProductID: p.ID,
			Action:    "CREATED",
		})

		// 3. Return HTTP 201 Created immediately to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
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

		if err := repo.UpdateProduct(r.Context(), id, p); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}
