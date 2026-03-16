package models

import "time"

const (
	OrderTypeMarket = "MARKET"
	OrderTypeLimit  = "LIMIT"
	OrderTypeStop   = "STOP"
)

var ValidOrderTypes = map[string]bool{
	OrderTypeMarket: true,
	OrderTypeLimit:  true,
	OrderTypeStop:   true,
}

type Order struct {
	ID         string `json:"id"`
	Asset      string `json:"asset"`
	Quantity   int    `json:"quantity"`
	OrderType  string `json:"order_type"`
	State      string `json:"state"`
	RetryCount int    `json:"retry_count"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type OrderEvent struct {
	OrderID       string `json:"order_id"`
	PreviousState string `json:"previous_state"`
	NewState      string `json:"new_state"`
	Timestamp     string `json:"timestamp"`
	Message       string `json:"message,omitempty"`
}

type CreateOrderRequest struct {
	Asset     string `json:"asset"`
	Quantity  int    `json:"quantity"`
	OrderType string `json:"order_type"`
}

func (r *CreateOrderRequest) Validate() error {
	if r.Asset == "" {
		return &ValidationError{Field: "asset", Message: "asset cannot be empty"}
	}
	if r.Quantity <= 0 {
		return &ValidationError{Field: "quantity", Message: "quantity must be greater than 0"}
	}
	if !ValidOrderTypes[r.OrderType] {
		return &ValidationError{Field: "order_type", Message: "order_type must be MARKET, LIMIT, or STOP"}
	}
	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

func NewOrder(id string, req *CreateOrderRequest) *Order {
	now := time.Now().UTC().Format(time.RFC3339)
	return &Order{
		ID:        id,
		Asset:     req.Asset,
		Quantity:  req.Quantity,
		OrderType: req.OrderType,
		State:     StateCreated,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
