package store

import "sync"

type Metrics struct {
	mu       sync.Mutex
	counters map[string]int
}

func NewMetrics() *Metrics {
	return &Metrics{
		counters: map[string]int{
			"total_orders":      0,
			"processing_orders": 0,
			"completed_orders":  0,
			"failed_orders":     0,
		},
	}
}

func (m *Metrics) Increment(key string, delta int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[key] += delta
}

func (m *Metrics) Decrement(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.counters[key] > 0 {
		m.counters[key]--
	}
}

func (m *Metrics) GetAll() map[string]int {
	m.mu.Lock()
	defer m.mu.Unlock()

	snapshot := make(map[string]int, len(m.counters))
	for k, v := range m.counters {
		snapshot[k] = v
	}
	return snapshot
}
