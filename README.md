Architectural Overview
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ HTTP/JSON
       ▼
┌─────────────┐
│  HTTP Server│ (net/http)
│  + Router   │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Handlers   │ (API Layer)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Services   │ (Business Logic)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Repository  │ (Data Access)
└──────┬──────┘
       │ database/sql
       ▼
┌─────────────┐
│  PostgreSQL │
└─────────────┘

## API Endpoints

### Authentication
- `POST /api/register` - Create new account
- `POST /api/login` - Login and get session token
- `POST /api/logout` - Logout and invalidate session
- `GET /api/me` - Get current user info

### Account Management
- `GET /api/account` - Get account details
- `GET /api/account/balance` - Get current balance
- `PATCH /api/account` - Update account details

### Transactions
- `POST /api/deposit` - Deposit money
- `POST /api/withdraw` - Withdraw money
- `POST /api/transfer` - Transfer to another account
- `GET /api/transactions` - Get transaction history (with pagination)
- `GET /api/transactions/:id` - Get specific transaction

### Health Check
- `GET /health` - Server health check

## Project Structure
```
bank-server/
├── main.go
├── go.mod
├── go.sum
├── config/
│   └── config.go          # Database config, server settings
├── db/
│   ├── db.go              # Database connection
│   └── migrations/
│       └── schema.sql     # Your schema
├── models/
│   ├── account.go
│   ├── transaction.go
│   └── session.go
├── repository/
│   ├── account_repo.go
│   └── transaction_repo.go
├── service/
│   ├── auth_service.go
│   └── transaction_service.go
├── handlers/
│   ├── auth_handler.go
│   ├── account_handler.go
│   └── transaction_handler.go
├── middleware/
│   ├── auth.go
│   └── logging.go
└── utils/
    ├── password.go        # bcrypt hashing
    ├── response.go        # JSON responses
    └── validation.go
