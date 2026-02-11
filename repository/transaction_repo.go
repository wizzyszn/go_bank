package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/wizzyszn/go_bank/db"
	"github.com/wizzyszn/go_bank/models"
)

type TransactionRepositoty struct {
	db *db.DB
}

func NewTransactionRepository(db *db.DB) *TransactionRepositoty {
	return &TransactionRepositoty{db: db}
}

func (t *TransactionRepositoty) Create(tx *sql.Tx, fromAccountID, toAccountID *int, amount float64, transactionType, description string) (*models.Transaction, error) {
	query := `
	INSERT INTO transactions(from_account_id,to_account_id,amount,type,description,status)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id,from_account_id,to_account_id,amount,type,description,status
	`
	transactions := &models.Transaction{}
	var err error

	if tx != nil {
		err = tx.QueryRow(query, fromAccountID, toAccountID, amount, transactionType, description, models.TransactionStatusCompleted).Scan(&transactions.FromAccountID, &transactions.ToAccountID, &transactions.Amount, &transactions.Type, &transactions.Description, &transactions.Status, &transactions.Status)
	} else {
		err = t.db.QueryRow(query, fromAccountID, toAccountID, amount, transactionType, description, models.TransactionStatusCompleted).Scan(&transactions.FromAccountID, &transactions.ToAccountID, &transactions.Amount, &transactions.Type, &transactions.Description, &transactions.Status, &transactions.Status)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return transactions, nil
}

func (r *TransactionRepositoty) GetByID(id int) (*models.Transaction, error) {
	query := `
	SELECT id,from_account_id,to_account_id,amount,type,description,status,created_at
	FROM transactions
	WHERE id = $1
	`
	transaction := &models.Transaction{}
	err := r.db.QueryRow(query, id).Scan(&transaction.ID, &transaction.FromAccountID, &transaction.ToAccountID, &transaction.Amount, &transaction.Type, &transaction.Description, &transaction.Status, &transaction.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("transactions not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

func (r *TransactionRepositoty) GetByAccountID(accountID, page, limit int) ([]*models.Transaction, int, error) {
	offset := (page - 1) * limit

	var totalCount int

	countQuery := `
	SELECT COUNT(*)
	FROM transactions
	WHERE from_account_id = $1 OR to_account_id = $1
	`

	err := r.db.QueryRow(countQuery, accountID).Scan(&totalCount)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	query := `
	SELECT id,from_account_id,to_account_id,amount,type,description,status,created_at
	FROM transactions
	WHERE from_account_id = $1 OR to_account_id = $1
	ORDER BY created_at DESC
	LIMIT $2
	OFFSET $3
	`
	rows, err := r.db.Query(query, accountID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	transactions := make([]*models.Transaction, 0)

	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(&transaction.ID, &transaction.FromAccountID, &transaction.ToAccountID, &transaction.Amount, &transaction.Type, &transaction.Description, &transaction.Status, &transaction.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, totalCount, nil
}

func (r *TransactionRepositoty) GetRecent(accountID, limit int) ([]*models.Transaction, error) {
	query := `
	SELECT id, from_account_id, to_account_id, amount, type, description, status, created_at
	FROM accounts
	WHERE from_account_id = $1 OR to_account_id = $2
	ORDER BY created_at DESC
	LIMIT $2
	`
	rows, err := r.db.Query(query, accountID, limit)

	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}
	defer rows.Close()

	transactions := make([]*models.Transaction, 0)

	for rows.Next() {
		transaction := &models.Transaction{}

		err := rows.Scan(&transaction.ID, transaction.FromAccountID, transaction.ToAccountID, transaction.Amount, transaction.Type, transaction.Status, transaction.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}
	return transactions, nil
}

func (r *TransactionRepositoty) GetByDateRange(accountID int, startDate, endDate time.Time) ([]*models.Transaction, error) {
	query := `
	SELECT id, from_account_id, to_account_id, amount, type, description, status, created_at
	FROM accounts
	WHERE (from_account_id = $1 OR to_account_id = $1)
	AND created_at >= $2
	AND created_at <= $3	
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, accountID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by date range: %w", err)
	}
	transactions := make([]*models.Transaction, 0)
	defer rows.Close()
	for rows.Next() {
		transaction := &models.Transaction{}

		err := rows.Scan(&transaction.ID, transaction.FromAccountID, transaction.ToAccountID, transaction.Amount, transaction.Type, transaction.Status, transaction.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("err iterating transactions: %w", err)
	}
	return transactions, nil
}

// GetTotalBalance
func (r *TransactionRepositoty) GetTotalBalance(accountID int) (float64, error) {
	var totalBalance float64

	query := `
	SELECT 
		COALESCE(SUM(CASE WHEN to_account_id = $1 THEN amount ELSE 0 END), 0) -
		COALESCE(SUM(CASE WHEN from_account_id = $1 THEN amount ELSE 0 END), 0)
	FROM transactions
	WHERE (from_account_id = $1 OR to_account_id = $1)
	AND status = $2
	`

	err := r.db.QueryRow(query, accountID, models.TransactionStatusCompleted).Scan(&totalBalance)
	if err != nil {
		return 0, fmt.Errorf("failed to get total balance: %w", err)
	}

	return totalBalance, nil
}
