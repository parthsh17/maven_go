package router

import (
	"maven/internal/handlers"
	"maven/internal/store"
	"maven/internal/worker"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewRouter(s *store.Store, m *store.Metrics, p *worker.Pool) http.Handler {
	orderHandler := handlers.NewOrderHandler(s, m, p)
	metricsHandler := handlers.NewMetricsHandler(m, p)

	mux := http.NewServeMux()

	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			orderHandler.CreateOrder(w, r)
		case http.MethodGet:
			orderHandler.ListOrders(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if len(path) > len("/orders/") && path[len(path)-len("/events"):] == "/events" {
			orderHandler.GetOrderEvents(w, r)
			return
		}

		orderHandler.GetOrder(w, r)
	})

	mux.HandleFunc("/metrics", metricsHandler.GetMetrics)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"maven"}`))
	})

	return corsMiddleware(mux)
}
