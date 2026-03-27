package models

import "fmt"

const (
	StateCreated   = "CREATED"
	StateValidated = "VALIDATED"
	StateQueued    = "QUEUED"
	StateExecuting = "EXECUTING"
	StateCompleted = "COMPLETED"
	StateFailed    = "FAILED"
	StateRetrying  = "RETRYING"
)

var ValidTransitions = map[string][]string{
	StateCreated:   {StateValidated},
	StateValidated: {StateQueued},
	StateQueued:    {StateExecuting},
	StateExecuting: {StateCompleted, StateFailed},
	StateFailed:    {StateRetrying},
	StateRetrying:  {StateQueued},
	StateCompleted: {},
}

type TransitionError struct {
	From string
	To   string
}

func (e *TransitionError) Error() string {
	return fmt.Sprintf("invalid state transition: %s → %s", e.From, e.To)
}

func CanTransition(from, to string) bool {
	allowed, exists := ValidTransitions[from]
	if !exists {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

func AllStates() []string {
	return []string{
		StateCreated, StateValidated, StateQueued,
		StateExecuting, StateCompleted, StateFailed, StateRetrying,
	}
}
