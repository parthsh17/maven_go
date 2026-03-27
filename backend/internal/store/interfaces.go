package store

import "maven/internal/models"

type UserStore interface {
	Register(id, email, password string) (*models.User, error)
	Authenticate(email, password string) (*models.User, error)
	GetUser(email string) (*models.User, error)
}

type OrderStore interface {
	AddOrder(order *models.Order) error
	GetOrder(id string) (*models.Order, error)
	GetAllOrders() []*models.Order
	UpdateState(id, newState, message string) error
	UpdateSlippage(id string, slippage float64) error
	IncrementRetry(id string)
	GetEvents(id string) ([]models.OrderEvent, error)
}
