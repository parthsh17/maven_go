package store

import "sync"

type Metrics struct {
	mu             sync.RWMutex
	counters       map[string]int
	successHistory []bool
	totalSlippage  float64
	slippageCount  int
}

func NewMetrics() *Metrics {
	return &Metrics{
		counters: map[string]int{
			"total_orders":      0,
			"processing_orders": 0,
			"completed_orders":  0,
			"failed_orders":     0,
		},
		successHistory: make([]bool, 0, 50),
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

func (m *Metrics) RecordResult(success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.successHistory) >= 50 {
		m.successHistory = m.successHistory[1:]
	}
	m.successHistory = append(m.successHistory, success)
}

func (m *Metrics) RecordSlippage(val float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalSlippage += val
	m.slippageCount++
}

func (m *Metrics) GetAll() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snapshot := make(map[string]interface{}, len(m.counters)+2)
	for k, v := range m.counters {
		snapshot[k] = v
	}

	// Calculate Moving Average Success Rate
	if len(m.successHistory) > 0 {
		successes := 0
		for _, s := range m.successHistory {
			if s {
				successes++
			}
		}
		snapshot["success_rate_ma"] = float64(successes) / float64(len(m.successHistory))
	} else {
		snapshot["success_rate_ma"] = 0.0
	}

	// Calculate Avg Slippage
	if m.slippageCount > 0 {
		snapshot["avg_slippage"] = m.totalSlippage / float64(m.slippageCount)
	} else {
		snapshot["avg_slippage"] = 0.0
	}

	return snapshot
}
