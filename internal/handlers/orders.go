package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"maven/internal/models"
	"maven/internal/store"
	"maven/internal/worker"
	"net/http"
	"strings"
)

func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type OrderHandler struct {
	store   *store.Store
	metrics *store.Metrics
	pool    *worker.Pool
}

func NewOrderHandler(s *store.Store, m *store.Metrics, p *worker.Pool) *OrderHandler {
	return &OrderHandler{store: s, metrics: m, pool: p}
}

func errorResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	id := newUUID()
	order := models.NewOrder(id, &req)

	if err := h.store.AddOrder(order); err != nil {
		errorResponse(w, http.StatusInternalServerError, "failed to store order: "+err.Error())
		return
	}

	h.metrics.Increment("total_orders", 1)

	if err := h.pool.Submit(order); err != nil {

		jsonResponse(w, http.StatusAccepted, map[string]interface{}{
			"order":   order,
			"warning": err.Error(),
		})
		return
	}

	jsonResponse(w, http.StatusCreated, order)
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	orders := h.store.GetAllOrders()
	if orders == nil {
		orders = []*models.Order{}
	}
	jsonResponse(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := extractID(r.URL.Path, "/orders/", "")
	if id == "" {
		errorResponse(w, http.StatusBadRequest, "missing order id")
		return
	}

	order, err := h.store.GetOrder(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, order)
}

func (h *OrderHandler) GetOrderEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := extractID(r.URL.Path, "/orders/", "/events")
	if id == "" {
		errorResponse(w, http.StatusBadRequest, "missing order id")
		return
	}

	events, err := h.store.GetEvents(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	if events == nil {
		events = []models.OrderEvent{}
	}
	jsonResponse(w, http.StatusOK, events)
}

func extractID(path, prefix, suffix string) string {
	path = strings.TrimPrefix(path, prefix)
	if suffix != "" {
		path = strings.TrimSuffix(path, suffix)
	}
	return strings.TrimSpace(path)
}
