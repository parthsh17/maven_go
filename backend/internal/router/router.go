package router

import (
	"log"
	"maven/internal/handlers"
	"maven/internal/store"
	"maven/internal/worker"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
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

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("[INFO] %s %s %d %v", r.Method, r.URL.Path, rw.status, time.Since(start))
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(5, 10)
		limiters[ip] = limiter
	}

	return limiter
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"too many requests","retry_after":"1s"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewRouter(s store.OrderStore, us store.UserStore, m *store.Metrics, p *worker.Pool) http.Handler {
	orderHandler := handlers.NewOrderHandler(s, m, p)
	metricsHandler := handlers.NewMetricsHandler(m, p)
	authHandler := handlers.NewAuthHandler(us)

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/signup", func(w http.ResponseWriter, r *http.Request) {
		authHandler.Signup(w, r)
	})

	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		authHandler.Login(w, r)
	})

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

	return recoveryMiddleware(loggingMiddleware(rateLimitMiddleware(corsMiddleware(mux))))
}
