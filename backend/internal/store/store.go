package store

import (
	"fmt"
	"maven/internal/models"
	"sync"
	"time"
)

type OrderStoreStruct struct {
	mu     sync.RWMutex
	orders map[string]*models.Order
	events map[string][]models.OrderEvent
}

func NewStore() OrderStore {
	return &OrderStoreStruct{
		orders: make(map[string]*models.Order),
		events: make(map[string][]models.OrderEvent),
	}
}

func (s *OrderStoreStruct) AddOrder(order *models.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.orders[order.ID]; exists {
		return fmt.Errorf("order %s already exists", order.ID)
	}
	s.orders[order.ID] = order
	return nil
}

func (s *OrderStoreStruct) GetOrder(id string) (*models.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.orders[id]
	if !exists {
		return nil, fmt.Errorf("order %s not found", id)
	}
	return order, nil
}

func (s *OrderStoreStruct) GetAllOrders() []*models.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*models.Order, 0, len(s.orders))
	for _, o := range s.orders {
		result = append(result, o)
	}
	return result
}

func (s *OrderStoreStruct) UpdateState(id, newState, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[id]
	if !exists {
		return fmt.Errorf("order %s not found", id)
	}

	if !models.CanTransition(order.State, newState) {
		return &models.TransitionError{From: order.State, To: newState}
	}

	prev := order.State
	order.State = newState
	order.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	event := models.OrderEvent{
		OrderID:       id,
		PreviousState: prev,
		NewState:      newState,
		Timestamp:     order.UpdatedAt,
		Message:       message,
	}
	s.events[id] = append(s.events[id], event)

	return nil
}

func (s *OrderStoreStruct) UpdateSlippage(id string, slippage float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[id]
	if !exists {
		return fmt.Errorf("order %s not found", id)
	}

	order.Slippage = slippage
	order.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (s *OrderStoreStruct) IncrementRetry(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if order, exists := s.orders[id]; exists {
		order.RetryCount++
		order.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
}

func (s *OrderStoreStruct) GetEvents(id string) ([]models.OrderEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, exists := s.events[id]
	if !exists {

		if _, ok := s.orders[id]; !ok {
			return nil, fmt.Errorf("order %s not found", id)
		}
		return []models.OrderEvent{}, nil
	}
	return events, nil
}
