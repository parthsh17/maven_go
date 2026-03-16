package handlers

import (
	"encoding/json"
	"maven/internal/store"
	"maven/internal/worker"
	"net/http"
)

type MetricsHandler struct {
	metrics *store.Metrics
	pool    *worker.Pool
}

func NewMetricsHandler(m *store.Metrics, p *worker.Pool) *MetricsHandler {
	return &MetricsHandler{metrics: m, pool: p}
}

func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	counters := h.metrics.GetAll()
	counters["worker_count"] = h.pool.WorkerCount()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(counters)
}
