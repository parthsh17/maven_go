# Maven: Order Lifecycle Management Platform

**Maven** is a high-throughput, concurrent order processing engine built with **Go** and **React**. It is designed to handle the full lifecycle of financial orders—from submission and validation to execution and auditing—with built-in resilience and real-time analytics.

---

## Key Features

### 1. Concurrent Worker Pool (Goroutines & Channels)
- **Concept**: Employs a fixed-size pool of worker Goroutines that consume orders from a shared buffered channel.
- **How it works**: When an order is submitted, it is placed into a `chan *models.Order`. Workers listen for incoming orders and process them in parallel. This decouples API request handling from the intensive work of order execution, allowing the system to scale efficiently.
- **Resilience**: Includes an **Automatic Retry Mechanism** that re-queues failed orders up to 3 times before marking them as permanently failed.

### 2. Advanced Middleware Engine
- **Structured Logging**: Every HTTP request is wrapped by a logging middleware that captures the method, URI, status code, and precise execution time (using `time.Since`). This provides critical visibility into API performance.
- **Panic Recovery**: A safety layer that uses `defer` and `recover()` to catch any unexpected runtime errors within handlers, preventing the entire server from crashing.
- **Rate Limiting**: Implements a token-bucket algorithm via `golang.org/x/time/rate`. It restricts users to 5 requests per second with a burst capacity of 10, protecting the system from brute-force attacks and resource exhaustion.
- **CORS Handler**: Manages cross-origin resource sharing to allow secure communication with the React frontend.

### 3. Database & Persistence (MongoDB)
- **State Persistence**: All orders, user accounts, and lifecycle events are persisted in **MongoDB**. This ensures data durability across server restarts.
- **Audit Logs**: Every state transition (e.g., `CREATED` → `EXECUTING` → `COMPLETED`) is recorded as a unique `OrderEvent` document, providing a full historical audit trail for every order.

### 4. Secure Authentication (Bcrypt)
- **Password Security**: Instead of storing plain-text passwords, the system uses **Bcrypt** hashing (`golang.org/x/crypto/bcrypt`).
- **How it works**: During signup, the password is salted and hashed before storage. During login, the provided password is compared against the stored hash using `CompareHashAndPassword`, ensuring that even if the database is compromised, user credentials remain secure.

### 5. Mathematical & Statistical Concepts
- **Slippage Calculation**: Simulates the real-world financial phenomenon where the execution price differs from the requested price.
  - **Implementation**: Uses a mathematical model `rand.Float64() * 0.01` to calculate a random slippage (0% to 1%) for each order, which is then recorded for quality analysis.
- **Moving Average Success Rate**: Tracks the reliability of the system over time.
  - **Implementation**: Uses a **Circular Buffer** (slice of size 50) to store the outcomes of the last 50 orders. The system calculates a rolling average of successful vs. failed orders, providing a more accurate performance metric than a simple lifetime total.
- **Regex-based Validation**: Uses the Regular Expression `^[A-Z]{3,5}/[A-Z]{3,5}$` to strictly validate currency pair formats (e.g., `BTC/USD`), ensuring data integrity at the entry point.

---

## Go Implementation Portfolio

This project serves as a showcase of advanced Go engineering patterns:
- **Interfaces**: Decoupled design using interfaces for `OrderStore`, `UserStore`, and `Executor`.
- **Sync Package**: Extensive use of `sync.RWMutex` to ensure thread-safety across concurrent metrics and in-memory lookups.
- **Context Handling**: Proper use of `context.Context` with timeouts for all MongoDB operations.

---

## Project Structure

```text
.
├── backend/            # Go Backend Service
│   ├── cmd/server/     # Entry point (main.go)
│   ├── internal/       # Core logic (Domain-Driven Design)
│   │   ├── handlers/   # HTTP handlers (REST API)
│   │   ├── store/      # Database (MongoDB) & Metrics logic
│   │   ├── worker/     # Concurrency / Worker Pool
│   │   └── models/     # Data structures & Validation (Regex)
│   └── tests/          # Robust unit and integration tests
└── frontend/           # React + Vite + TypeScript Dashboard
```

---

## Getting Started

### Backend
1. `cd backend`
2. `go run cmd/server/main.go`
- API runs on `http://localhost:8080`

### Frontend
1. `cd frontend`
2. `npm install`
3. `npm run dev`
- Dashboard runs on `http://localhost:5173`
