# 🏦 GoBank

A production-ready banking REST API built from scratch in Go using only the standard library (`net/http`) and PostgreSQL. Features session-based authentication, transactional money operations, middleware chaining, and graceful shutdown.

---

## ✨ Features

- **Authentication** — Register, login, logout with session-based token auth
- **Account Management** — View & update account details, check balances
- **Transactions** — Deposits, withdrawals, and account-to-account transfers with database transactions
- **Middleware Pipeline** — Composable middleware chain with logging, CORS, rate limiting, and authentication
- **Rate Limiting** — Token bucket algorithm with per-IP tracking and `Retry-After` headers
- **CORS** — Environment-aware CORS (permissive in development, locked-down in production)
- **Health Checks** — `/health`, `/ready`, and `/live` endpoints (Kubernetes-compatible)
- **Session Cleanup** — Background goroutine purges expired sessions every hour
- **Graceful Shutdown** — Signal-based shutdown with a 30-second drain period
- **Input Validation** — Request validation with structured error responses
- **Password Security** — bcrypt hashing for all stored passwords

---

## 🏗 Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ HTTP/JSON
       ▼
┌─────────────┐
│  HTTP Server│  net/http + middleware chain
│  + Router   │  (Logger → CORS → RateLimit → Auth)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Handlers   │  Request parsing, response writing
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Services   │  Business logic, validation, transactions
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Repositories│  SQL queries, data access
└──────┬──────┘
       │ database/sql + lib/pq
       ▼
┌─────────────┐
│ PostgreSQL  │
└─────────────┘
```

---

## 📁 Project Structure

```
go_bank/
├── main.go                          # Entrypoint: wiring, routing, server lifecycle
├── config/
│   ├── config.go                    # Env-based configuration (database, server, security)
│   └── config_test.go
├── db/
│   ├── db.go                        # Connection pool, health checks, stats
│   ├── migrate.go                   # Schema migration runner
│   ├── transaction.go               # DB transaction helper (Begin/Commit/Rollback)
│   └── migrations/
│       └── schema.sql               # Full schema: accounts, transactions, sessions
├── models/
│   ├── account.go                   # Account model, request/response types
│   ├── transaction.go               # Transaction model, request/response types
│   ├── session.go                   # Session model
│   └── response.go                  # Generic API response wrapper
├── repository/
│   ├── account_repo.go              # Account CRUD operations
│   ├── session_repo.go              # Session CRUD + cleanup
│   ├── transaction_repo.go          # Transaction queries + pagination
│   └── transaction_repo_test.go
├── service/
│   ├── auth_service.go              # Registration, login, logout, session mgmt
│   └── transaction_service.go       # Deposit, withdraw, transfer, balance
├── handlers/
│   ├── auth_handler.go              # POST /register, /login, /logout; GET /me
│   ├── account_handler.go           # GET/PATCH /account, GET /account/balance
│   ├── transaction_handler.go       # POST /deposit, /withdraw, /transfer; GET /transactions
│   └── health_handler.go            # GET /health, /ready, /live
├── middleware/
│   ├── chain.go                     # Middleware chaining utility
│   ├── auth.go                      # Session-based authentication middleware
│   ├── cors.go                      # CORS (dev + production configs)
│   ├── ratelimit.go                 # Token bucket rate limiter
│   └── logging.go                   # Request/response logger
├── utils/
│   ├── password.go                  # bcrypt hash + compare
│   ├── response.go                  # JSON response helpers (success, error, etc.)
│   ├── session.go                   # Session token generation
│   ├── validation.go                # Input validation + ValidationError type
│   └── utils_test.go
├── scripts/
│   ├── setup_db.sh                  # Interactive PostgreSQL database setup
│   └── run_migration.sh             # Migration runner script
├── .env                             # Environment variables (not for production)
├── .gitignore
├── go.mod
└── go.sum
```

---

## 🚀 Getting Started

### Prerequisites

- **Go** 1.25+
- **PostgreSQL** 12+

### 1. Clone the repository

```bash
git clone https://github.com/wizzyszn/go_bank.git
cd go_bank
```

### 2. Set up the database

Use the interactive setup script:

```bash
chmod +x scripts/setup_db.sh
./scripts/setup_db.sh
```

Then run the schema migration:

```bash
chmod +x scripts/run_migration.sh
./scripts/run_migration.sh
```

### 3. Configure environment variables

Copy and edit the `.env` file:

```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=bankdb
DB_SSLMODE=disable

# Security
SESSION_SECRET=change-this-to-a-random-secret-in-production
SESSION_DURATION_HOURS=24
```

### 4. Run the server

```bash
go run main.go
```

The server starts on `http://localhost:8080`.

---

## 📡 API Reference

All endpoints return JSON. Protected endpoints require an `Authorization: Bearer <session_token>` header.

### Health

| Method | Endpoint  | Description                     |
| ------ | --------- | ------------------------------- |
| GET    | `/health` | Server + database health status |
| GET    | `/ready`  | Readiness probe (DB ping)       |
| GET    | `/live`   | Liveness probe (always 200)     |

### Authentication (Public)

| Method | Endpoint        | Description                  |
| ------ | --------------- | ---------------------------- |
| POST   | `/api/register` | Create a new account         |
| POST   | `/api/login`    | Login, returns session token |

### Authentication (Protected)

| Method | Endpoint      | Description                      |
| ------ | ------------- | -------------------------------- |
| POST   | `/api/logout` | Invalidate current session       |
| GET    | `/api/me`     | Get authenticated user's profile |

### Account Management (Protected)

| Method | Endpoint               | Description                     |
| ------ | ---------------------- | ------------------------------- |
| GET    | `/api/account`         | Get account details             |
| PATCH  | `/api/account`         | Update account (name, password) |
| GET    | `/api/account/balance` | Get current balance             |

### Transactions (Protected)

| Method | Endpoint            | Description                                    |
| ------ | ------------------- | ---------------------------------------------- |
| POST   | `/api/deposit`      | Deposit funds                                  |
| POST   | `/api/withdraw`     | Withdraw funds                                 |
| POST   | `/api/transfer`     | Transfer funds to another account              |
| GET    | `/api/transactions` | List transactions (paginated: `?page=&limit=`) |

---

## 📝 Example Requests

### Register

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "password": "securepassword123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### Deposit

```bash
curl -X POST http://localhost:8080/api/deposit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <session_token>" \
  -d '{
    "amount": 500.00,
    "description": "Initial deposit"
  }'
```

### Transfer

```bash
curl -X POST http://localhost:8080/api/transfer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <session_token>" \
  -d '{
    "to_account_id": 2,
    "amount": 100.00,
    "description": "Payment to Jane"
  }'
```

---

## 🗄 Database Schema

Three core tables with proper constraints, indexes, and triggers:

- **`accounts`** — User accounts with email, hashed password, balance (non-negative constraint), currency, and status
- **`transactions`** — Financial records with foreign keys to sender/receiver, amount (positive constraint), type, and status
- **`sessions`** — Token-based sessions with expiration; auto-cleaned by background job and a PostgreSQL function

---

## 🔒 Security

- Passwords hashed with **bcrypt**
- Session tokens for stateful authentication
- Rate limiting with **token bucket** algorithm (30 req/s per IP, burst of 100)
- CORS configured per environment
- Input validation on all endpoints
- Sensitive fields (e.g. `password_hash`) stripped from API responses

---

## 🧪 Testing

```bash
go test ./...
```

---

## 📦 Dependencies

| Package                    | Purpose                 |
| -------------------------- | ----------------------- |
| `github.com/lib/pq`        | PostgreSQL driver       |
| `github.com/joho/godotenv` | `.env` file loading     |
| `golang.org/x/crypto`      | bcrypt password hashing |

---

## 📄 License

MIT
