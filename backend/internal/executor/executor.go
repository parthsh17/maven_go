package executor

import (
	"fmt"
	"math/rand"
)

type Executor interface {
	Execute(orderID string) (float64, error)
}

type BasicExecutor struct {
	successRate float64
}

func NewBasicExecutor(successRate float64) *BasicExecutor {
	return &BasicExecutor{successRate: successRate}
}

func (e *BasicExecutor) Execute(orderID string) (float64, error) {
	slippage := rand.Float64() * 0.01 // Simulating up to 1% slippage

	if rand.Float64() < e.successRate {
		return slippage, nil
	}
	return 0, fmt.Errorf("execution failed for order %s: market conditions unfavorable", orderID)
}
