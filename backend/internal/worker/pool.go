package worker

import (
	"fmt"
	"log"
	"maven/internal/executor"
	"maven/internal/models"
	"maven/internal/store"
	"sync"
)

const MaxRetries = 3

type Pool struct {
	workerCount int
	orderCh     chan *models.Order
	store       store.OrderStore
	metrics     *store.Metrics
	executor    executor.Executor
	wg          sync.WaitGroup
	once        sync.Once
}

func NewPool(workerCount int, bufferSize int, s store.OrderStore, m *store.Metrics, exec executor.Executor) *Pool {
	return &Pool{
		workerCount: workerCount,
		orderCh:     make(chan *models.Order, bufferSize),
		store:       s,
		metrics:     m,
		executor:    exec,
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.runWorker(i)
	}
	log.Printf("Worker pool started with %d workers", p.workerCount)
}

func (p *Pool) Submit(order *models.Order) error {
	select {
	case p.orderCh <- order:
		return nil
	default:
		return fmt.Errorf("worker pool queue is full, try again later")
	}
}

func (p *Pool) Stop() {
	p.once.Do(func() {
		close(p.orderCh)
	})
	p.wg.Wait()
	log.Println("Worker pool stopped")
}

func (p *Pool) WorkerCount() int {
	return p.workerCount
}

func (p *Pool) runWorker(id int) {
	defer p.wg.Done()

	for order := range p.orderCh {
		log.Printf("[Worker %d] Processing order %s (state: %s)", id, order.ID, order.State)
		p.processOrder(order)
	}

	log.Printf("[Worker %d] Shutting down", id)
}

func (p *Pool) processOrder(order *models.Order) {
	id := order.ID

	current, err := p.store.GetOrder(id)
	if err != nil {
		log.Printf("Error fetching order %s: %v", id, err)
		return
	}

	if current.State == models.StateCreated {
		if err := p.store.UpdateState(id, models.StateValidated, "order validated"); err != nil {
			log.Printf("Error validating order %s: %v", id, err)
			return
		}
		if err := p.store.UpdateState(id, models.StateQueued, "order queued for execution"); err != nil {
			log.Printf("Error queuing order %s: %v", id, err)
			return
		}
		p.metrics.Increment("processing_orders", 1)
	}

	if err := p.store.UpdateState(id, models.StateExecuting, "execution started"); err != nil {
		log.Printf("Error executing order %s: %v", id, err)
		if current.State == models.StateCreated {
			p.metrics.Decrement("processing_orders")
		}
		return
	}

	slippage, execErr := p.executor.Execute(id)

	if execErr == nil {
		if err := p.store.UpdateState(id, models.StateCompleted, "execution successful"); err != nil {
			log.Printf("Error completing order %s: %v", id, err)
		}
		if err := p.store.UpdateSlippage(id, slippage); err != nil {
			log.Printf("Error updating slippage for order %s: %v", id, err)
		}
		p.metrics.Decrement("processing_orders")
		p.metrics.Increment("completed_orders", 1)
		p.metrics.RecordResult(true)
		p.metrics.RecordSlippage(slippage)
		log.Printf("Order %s COMPLETED with %.4f%% slippage", id, slippage*100)
		return
	}

	p.metrics.RecordResult(false)
	if err := p.store.UpdateState(id, models.StateFailed, execErr.Error()); err != nil {
		log.Printf("Error failing order %s: %v", id, err)
		p.metrics.Decrement("processing_orders")
		return
	}

	p.store.IncrementRetry(id)
	retryOrder, err := p.store.GetOrder(id)
	if err != nil {
		log.Printf("Error fetching order %s for retry: %v", id, err)
		p.metrics.Decrement("processing_orders")
		p.metrics.Increment("failed_orders", 1)
		return
	}

	if retryOrder.RetryCount <= MaxRetries {

		if err := p.store.UpdateState(id, models.StateRetrying,
			fmt.Sprintf("retry attempt %d", retryOrder.RetryCount)); err != nil {
			log.Printf("Error setting retry state %s: %v", id, err)
		} else {
			if err := p.store.UpdateState(id, models.StateQueued, "re-queued after retry"); err != nil {
				log.Printf("Error re-queuing order %s: %v", id, err)
			} else {
				log.Printf("Order %s re-queued for retry #%d", id, retryOrder.RetryCount)

				_ = p.Submit(retryOrder)
				return
			}
		}
	}

	p.metrics.Decrement("processing_orders")
	p.metrics.Increment("failed_orders", 1)
	log.Printf("Order %s permanently FAILED after %d retries", id, retryOrder.RetryCount)
}
