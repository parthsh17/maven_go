package executor

import (
	"fmt"
	"math/rand"
)

type Executor interface {
	Execute(orderID string) error
}

type BasicExecutor struct {
	successRate float64
}

func NewBasicExecutor(successRate float64) *BasicExecutor {
	return &BasicExecutor{successRate: successRate}
}

func (e *BasicExecutor) Execute(orderID string) error {

	if rand.Float64() < e.successRate {
		return nil
	}
	return fmt.Errorf("execution failed for order %s: market conditions unfavorable", orderID)
}
