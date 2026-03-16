package main

import (
	"context"
	"log"
	"maven/internal/executor"
	"maven/internal/router"
	"maven/internal/store"
	"maven/internal/worker"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	addr        = ":8080"
	workerCount = 5
	bufferSize  = 100
	successRate = 0.70
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Maven Order Lifecycle Platform...")

	s := store.NewStore()
	us := store.NewUserStore()
	m := store.NewMetrics()

	exec := executor.NewBasicExecutor(successRate)

	pool := worker.NewPool(workerCount, bufferSize, s, m, exec)
	pool.Start()
	log.Printf("[CONCURRENCY] Worker pool started with %d concurrent goroutines", workerCount)

	handler := router.NewRouter(s, us, m, pool)

	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("Maven server running on %s", addr)
		errCh <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Fatalf("Server error: %v", err)
	case sig := <-quit:
		log.Printf("Received signal %v, shutting down...", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	pool.Stop()
	log.Println("Maven shut down cleanly.")
}
