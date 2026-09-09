package main

import (
	"Crochet/internal/handler"
	"Crochet/internal/middleware"
	"Crochet/internal/queue"
	"Crochet/internal/store"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func applyMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func main() {
	sqStore, err := store.NewSQLiteStore("crochet.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	taskQ := queue.NewTaskQueue(100, 2)

	mux := http.NewServeMux()

	// --- 1. Public Routes ---
	mux.HandleFunc("GET /health", handler.HandleHealth())
	mux.HandleFunc("GET /products", handler.HandleGetProducts(sqStore))

	// --- 2. Protected Routes (Require Authentication) ---
	mux.Handle("POST /products", middleware.RequireAuth(handler.HandleAddProduct(sqStore, taskQ)))
	mux.Handle("DELETE /products/{id}", middleware.RequireAuth(handler.HandleDeleteProduct(sqStore)))
	mux.Handle("PUT /products/{id}", middleware.RequireAuth(handler.HandleUpdateProduct(sqStore)))

	// Intentional crash route to test recovery middleware
	mux.HandleFunc("GET /panic-test", func(w http.ResponseWriter, r *http.Request) {
		var nilPointer *string
		_ = *nilPointer // Triggers an intentional panic!
	})

	// Wrap the entire Mux router with global middleware (Recoverer & Logger)
	wrappedMux := applyMiddleware(mux,
		middleware.Recoverer, // 1. Catches crashes across all routes
		middleware.Logger,    // 2. Logs execution time & status code
	)

	// 1. Configure custom http.Server instance
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080" // Default fallback if not set
	} else if !strings.HasPrefix(port, ":") {
		port = ":" + port // Guarantees leading colon (e.g., "8080" -> ":8080")
	}

	server := &http.Server{
		Addr:    port,
		Handler: wrappedMux,
	}

	// 2. Start HTTP server in a background goroutine so it doesn't block main
	go func() {
		log.Printf("Server running on http://localhost:%s...\n", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// 3. Create a channel to catch OS shutdown signals (Ctrl+C, SIGTERM)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// 4. BLOCK HERE until an OS signal is received!
	<-stopChan
	log.Println("\nShutdown signal received. Starting graceful shutdown...")

	// 5. Create a 10-second deadline context for HTTP server to finish active web requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 6. Stop accepting new HTTP requests and wait for in-flight requests to finish
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	} else {
		log.Println("HTTP server stopped gracefully.")
	}

	// 7. Drain and stop background worker queue
	taskQ.Stop()

	// 8. Close database connections cleanly
	if err := sqStore.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		log.Println("Database connection closed.")
	}

	log.Println("Application exited cleanly.")
}
