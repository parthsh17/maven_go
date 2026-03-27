package main

import (
	"context"
	"log"
	"maven/internal/config"
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

func main() {
	cfg := config.Load()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Maven")

	mongoClient, err := store.NewMongoClient(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer mongoClient.Disconnect()

	s := store.NewMongoOrderStore(mongoClient.DB)
	us := store.NewMongoUserStore(mongoClient.DB)
	m := store.NewMetrics()

	exec := executor.NewBasicExecutor(cfg.SuccessRate)

	pool := worker.NewPool(cfg.WorkerCount, cfg.BufferSize, s, m, exec)
	pool.Start()
	log.Printf("[CONCURRENCY] Worker pool started with %d concurrent goroutines", cfg.WorkerCount)

	handler := router.NewRouter(s, us, m, pool)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("Maven server running on %s", cfg.Addr)
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
	log.Println("Maven shut down")
}
