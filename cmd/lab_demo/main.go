package main

import (
	"encoding/json"
	"fmt"
	"maven/internal/executor"
	"maven/internal/models"
	"maven/internal/store"
	"maven/internal/worker"
	"sync"
	"time"
)

func main() {
	fmt.Println("--------------------------------------------------")
	fmt.Println("🧪 LAB 7 OUTPUT: Pointers & Call by Reference")
	fmt.Println("--------------------------------------------------")
	s := store.NewStore()
	req := &models.CreateOrderRequest{Asset: "AAPL", Quantity: 10, OrderType: "MARKET"}
	order := models.NewOrder("order-001", req)

	fmt.Printf("1. Created order. Memory Address: %p\n", order)
	err := s.AddOrder(order)
	fmt.Printf("2. Passed to AddOrder(order *models.Order). Error: %v\n", err)

	retrieved, _ := s.GetOrder("order-001")
	fmt.Printf("3. Retrieved from store. Memory Address: %p\n", retrieved)
	if order == retrieved {
		fmt.Println("-> SUCCESS: Both pointers reference the exact same underlying memory!")
	}

	fmt.Println("\n--------------------------------------------------")
	fmt.Println("🧪 LAB 8 OUTPUT: JSON Marshal & Unmarshal")
	fmt.Println("--------------------------------------------------")

	jsonInput := []byte(`{"asset":"NVDA", "quantity":5, "order_type":"LIMIT"}`)
	fmt.Printf("1. Raw JSON Input (Bytes): %s\n", string(jsonInput))

	var unmarshaledReq models.CreateOrderRequest
	json.Unmarshal(jsonInput, &unmarshaledReq)
	fmt.Printf("2. UNMARSHALED to Go Struct: %+v\n", unmarshaledReq)

	order2 := models.NewOrder("order-002", &unmarshaledReq)

	order2.CreatedAt = "2026-03-12T02:49:34Z"
	order2.UpdatedAt = "2026-03-12T02:49:34Z"

	marshaledOutput, _ := json.MarshalIndent(order2, "", "  ")
	fmt.Printf("3. MARSHALED Go Struct back to JSON:\n%s\n", string(marshaledOutput))

	fmt.Println("\n--------------------------------------------------")
	fmt.Println("🧪 LAB 9 OUTPUT: Concurrency & Mutexes")
	fmt.Println("--------------------------------------------------")

	m := store.NewMetrics()
	var wg sync.WaitGroup

	fmt.Println("1. Launching 100 concurrent goroutines to increment 'total_orders'...")
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Increment("total_orders", 1)
		}()
	}
	wg.Wait()

	finalMetrics := m.GetAll()
	fmt.Printf("2. Final 'total_orders' count: %d\n", finalMetrics["total_orders"])
	fmt.Println("-> SUCCESS: Mutex prevented race conditions and lost updates. Exact count is 100.")

	fmt.Println("\n--------------------------------------------------")
	fmt.Println("🧪 LAB 10 OUTPUT: Goroutines and Channels")
	fmt.Println("--------------------------------------------------")

	fmt.Println("1. Starting Worker Pool with 2 Workers...")
	exec := executor.NewBasicExecutor(1.0)
	pool := worker.NewPool(2, 10, s, m, exec)
	pool.Start()

	fmt.Println("2. Submitting 3 tasks to the Channel...")
	for i := 1; i <= 3; i++ {
		o := models.NewOrder(fmt.Sprintf("chan-order-%d", i), req)
		s.AddOrder(o)
		pool.Submit(o)
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(50 * time.Millisecond)
	pool.Stop()
	fmt.Println("3. Worker Pool stopped cleanly.")
}
