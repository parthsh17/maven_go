package tests

import (
	"fmt"
	"maven/internal/models"
	"maven/internal/store"
	"testing"
)

func TestNewOrder_Fields(t *testing.T) {
	req := &models.CreateOrderRequest{
		Asset:     "AAPL",
		Quantity:  10,
		OrderType: models.OrderTypeMarket,
	}
	order := models.NewOrder("order-001", req)

	if order.ID != "order-001" {
		t.Errorf("expected ID 'order-001', got %q", order.ID)
	}
	if order.Asset != "AAPL" {
		t.Errorf("expected Asset 'AAPL', got %q", order.Asset)
	}
	if order.Quantity != 10 {
		t.Errorf("expected Quantity 10, got %d", order.Quantity)
	}
	if order.State != models.StateCreated {
		t.Errorf("expected initial State %q, got %q", models.StateCreated, order.State)
	}
	if order.RetryCount != 0 {
		t.Errorf("expected RetryCount 0, got %d", order.RetryCount)
	}
}

func TestValidate_EmptyAsset(t *testing.T) {
	req := &models.CreateOrderRequest{Asset: "", Quantity: 5, OrderType: "MARKET"}
	if err := req.Validate(); err == nil {
		t.Error("expected error for empty asset, got nil")
	}
}

func TestValidate_ZeroQuantity(t *testing.T) {
	req := &models.CreateOrderRequest{Asset: "TSLA", Quantity: 0, OrderType: "MARKET"}
	if err := req.Validate(); err == nil {
		t.Error("expected error for zero quantity, got nil")
	}
}

func TestValidate_NegativeQuantity(t *testing.T) {
	req := &models.CreateOrderRequest{Asset: "TSLA", Quantity: -5, OrderType: "MARKET"}
	if err := req.Validate(); err == nil {
		t.Error("expected error for negative quantity, got nil")
	}
}

func TestValidate_InvalidOrderType(t *testing.T) {
	req := &models.CreateOrderRequest{Asset: "TSLA", Quantity: 5, OrderType: "INVALID"}
	if err := req.Validate(); err == nil {
		t.Error("expected error for invalid order type, got nil")
	}
}

func TestValidate_ValidRequest(t *testing.T) {
	for _, ot := range []string{"MARKET", "LIMIT", "STOP"} {
		req := &models.CreateOrderRequest{Asset: "GOOG", Quantity: 1, OrderType: ot}
		if err := req.Validate(); err != nil {
			t.Errorf("expected nil error for type %q, got %v", ot, err)
		}
	}
}

func TestCanTransition_ValidPaths(t *testing.T) {
	validPaths := [][2]string{
		{models.StateCreated, models.StateValidated},
		{models.StateValidated, models.StateQueued},
		{models.StateQueued, models.StateExecuting},
		{models.StateExecuting, models.StateCompleted},
		{models.StateExecuting, models.StateFailed},
		{models.StateFailed, models.StateRetrying},
		{models.StateRetrying, models.StateQueued},
	}
	for _, path := range validPaths {
		if !models.CanTransition(path[0], path[1]) {
			t.Errorf("expected %s → %s to be valid", path[0], path[1])
		}
	}
}

func TestCanTransition_InvalidPaths(t *testing.T) {
	invalidPaths := [][2]string{
		{models.StateCreated, models.StateCompleted},
		{models.StateCompleted, models.StateCreated},
		{models.StateValidated, models.StateExecuting},
		{models.StateCompleted, models.StateQueued},
	}
	for _, path := range invalidPaths {
		if models.CanTransition(path[0], path[1]) {
			t.Errorf("expected %s → %s to be INVALID, but got valid", path[0], path[1])
		}
	}
}

func TestStore_AddAndGetOrder(t *testing.T) {
	s := store.NewStore()
	req := &models.CreateOrderRequest{Asset: "NVDA", Quantity: 3, OrderType: "LIMIT"}
	order := models.NewOrder("test-id-1", req)

	if err := s.AddOrder(order); err != nil {
		t.Fatalf("AddOrder failed: %v", err)
	}

	got, err := s.GetOrder("test-id-1")
	if err != nil {
		t.Fatalf("GetOrder failed: %v", err)
	}

	if got.ID != "test-id-1" {
		t.Errorf("expected ID 'test-id-1', got %q", got.ID)
	}
	if got.Asset != "NVDA" {
		t.Errorf("expected Asset 'NVDA', got %q", got.Asset)
	}
}

func TestStore_DuplicateOrder(t *testing.T) {
	s := store.NewStore()
	req := &models.CreateOrderRequest{Asset: "MSFT", Quantity: 2, OrderType: "STOP"}
	order := models.NewOrder("dup-id", req)

	_ = s.AddOrder(order)
	if err := s.AddOrder(order); err == nil {
		t.Error("expected duplicate order error, got nil")
	}
}

func TestStore_OrderNotFound(t *testing.T) {
	s := store.NewStore()
	_, err := s.GetOrder("nonexistent")
	if err == nil {
		t.Error("expected not-found error, got nil")
	}
}

func TestStore_UpdateState_Valid(t *testing.T) {
	s := store.NewStore()
	req := &models.CreateOrderRequest{Asset: "AMZN", Quantity: 5, OrderType: "MARKET"}
	order := models.NewOrder("state-id", req)
	_ = s.AddOrder(order)

	if err := s.UpdateState("state-id", models.StateValidated, "validated"); err != nil {
		t.Fatalf("UpdateState to VALIDATED failed: %v", err)
	}
	got, _ := s.GetOrder("state-id")
	if got.State != models.StateValidated {
		t.Errorf("expected VALIDATED, got %q", got.State)
	}
}

func TestStore_UpdateState_Invalid(t *testing.T) {
	s := store.NewStore()
	req := &models.CreateOrderRequest{Asset: "AMZN", Quantity: 5, OrderType: "MARKET"}
	order := models.NewOrder("bad-transition", req)
	_ = s.AddOrder(order)

	if err := s.UpdateState("bad-transition", models.StateCompleted, "skip"); err == nil {
		t.Error("expected transition error, got nil")
	}
}

func TestStore_GetAllOrders(t *testing.T) {
	s := store.NewStore()
	for i := 0; i < 3; i++ {
		req := &models.CreateOrderRequest{Asset: "TEST", Quantity: 1, OrderType: "MARKET"}
		o := models.NewOrder(fmt.Sprintf("all-orders-test-%d", i), req)
		_ = s.AddOrder(o)
	}
	all := s.GetAllOrders()
	if len(all) != 3 {
		t.Errorf("expected 3 orders, got %d", len(all))
	}
}

func TestStore_EventLogging(t *testing.T) {
	s := store.NewStore()
	req := &models.CreateOrderRequest{Asset: "IBM", Quantity: 7, OrderType: "LIMIT"}
	order := models.NewOrder("event-id", req)
	_ = s.AddOrder(order)

	_ = s.UpdateState("event-id", models.StateValidated, "step 1")
	_ = s.UpdateState("event-id", models.StateQueued, "step 2")

	events, err := s.GetEvents("event-id")
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestMetrics_IncrementAndGet(t *testing.T) {
	m := store.NewMetrics()
	m.Increment("total_orders", 1)
	m.Increment("total_orders", 1)

	all := m.GetAll()
	if all["total_orders"] != 2 {
		t.Errorf("expected total_orders=2, got %d", all["total_orders"])
	}
}

func TestMetrics_DecrementFloor(t *testing.T) {
	m := store.NewMetrics()
	m.Decrement("processing_orders")
	all := m.GetAll()
	if all["processing_orders"] < 0 {
		t.Errorf("counter should not go negative, got %d", all["processing_orders"])
	}
}
