package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// Logger logs incoming requests, HTTP methods, latency, and status codes.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		rw := newResponseWriter(w)

		// Execute the next handler in the chain
		next.ServeHTTP(rw, r)

		// Code below runs AFTER the handler finishes!
		log.Printf("[%s] %s %s | Status: %d | Duration: %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			rw.statusCode,
			time.Since(start),
		)
	})
}

// Recoverer catches any unexpected panic, prevents the server from crashing,
// and returns a clean 500 Internal Server Error response to the client.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic message and stack trace
				log.Printf("[PANIC RECOVERED] %v\n%s", err, string(debug.Stack()))

				// Send clean 500 JSON response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, `{"error":"Internal Server Error"}`)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// --- Helper to capture HTTP status codes ---

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
