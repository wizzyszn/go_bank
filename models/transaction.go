package models

import "time"

// Transaction represents a financial transaction

type Transaction struct {
	ID            int       `json:"id" db:"id"`
	FromAccountID int       `json:"from_account_id" db:"from_account_id"`
	ToAccountID   int       `json:"to_account_id" db:"to_account_id"`
	Amount        float64   `json:"amount" db:"amount"`
	Type          string    `json:"type" db:"type"`
	Description   string    `json:"description" db:"description"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// DepositRequest represents a deposit request

type DepositRequest struct {
	Amount      float64 `json:"amount"`
	Description string `json:"description"`
}

// WithdrawRequest represents a withdrawal request

type WitdrawRequest struct {
	Amount      float64 `json:"amount"`
	Description string `json:"description"`
}

// TransferRequest represents a transfer request
type TransferRequest struct {
	ToAccountID int     `json:"to_account_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// TransactionResponse is what we return to the client
type TransactionResponse struct {
	ID            int       `json:"id"`
	FromAccountID *int      `json:"from_account_id,omitempty"`
	ToAccountID   *int      `json:"to_account_id,omitempty"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// ToResponse converts Transaction to TransactionResponse
func (t *Transaction) ToResponse() *TransactionResponse {
	return &TransactionResponse{
		ID:            t.ID,
		FromAccountID: &t.FromAccountID,
		ToAccountID:   &t.ToAccountID,
		Amount:        t.Amount,
		Type:          t.Type,
		Description:   t.Description,
		Status:        t.Status,
		CreatedAt:     t.CreatedAt,
	}
}

// Transaction types
const (
	TransactionTypeDeposit  string = "deposit"
	TransactionTypeTransfer string = "transfer"
	TransactionTypeWithdraw string = "withdraw"
)

// Transaction statuses

const (
	TransactionStatusPending   = "pending"
	TransactionStatusCompleted = "completed"
	TransactionStatusFailed    = "failed"
)
