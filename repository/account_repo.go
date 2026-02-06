package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/wizzyszn/go_bank/db"
	"github.com/wizzyszn/go_bank/models"
)

type AccountRepository struct {
	db *db.DB
}

func NewAccountRepository(database *db.DB) *AccountRepository {
	return &AccountRepository{
		db: database,
	}
}

func (r *AccountRepository) Create(email, passwordHash, firstName, lastName string) (*models.Account, error) {

	query := `
	INSERT INTO accounts (email,password_hash,first_name,last_name,balance,currency,status)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING id,email,first_name,last_name,balance,currency,status,created_at,updated_at
	`
	account := &models.Account{}

	err := r.db.QueryRow(query, email, passwordHash, firstName, lastName, 0.00, "USD", models.AccountStatusActice).Scan(&account.ID, &account.Email, &account.FirstName, &account.LastName, &account.Balance, &account.Currency, &account.Status, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (r *AccountRepository) GetByID(id int) (*models.Account, error) {
	query := `
	SELECT id,email,first_name,last_name,balance,currency,status,created_at,updated_at
	FROM accounts
	WHERE id = $1
	`
	account := &models.Account{}
	err := r.db.QueryRow(query, id).Scan(&account.ID, &account.Email, &account.FirstName, &account.LastName, &account.Balance, &account.Currency, &account.Status, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get an account: %w", err)
	}
	return account, nil
}

func (r *AccountRepository) GeyByEmail(email string) (*models.Account, error) {
	query := `
	SELECT id,email,first_name,last_name,balance,currency,status,created_at,updated_at
	FROM accounts
	WHERE email = $1
	`
	account := &models.Account{}
	err := r.db.QueryRow(query, email).Scan(&account.ID, &account.Email, &account.FirstName, &account.LastName, &account.Balance, &account.Currency, &account.Status, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get an account: %w", err)
	}

	return account, nil
}

func (r *AccountRepository) Update(id int, firstName, lastName string) error {
	query := `
	UPDATE accounts
	SET first_name = $1 ,last_name = $2, updated_at = $3
	WHERE id = $4
	`
	result, err := r.db.Exec(query, firstName, lastName, time.Now(), id)

	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

func (r *AccountRepository) UpdateBalance(tx *sql.Tx, accountID int, newBalace float64) error {

	query := `
	UPDATE accounts
	SET balance = $1 , updated_at = $2
	WHERE id = $3
	`
	var err error
	var result sql.Result

	if tx != nil {
		result, err = tx.Exec(query, newBalace, time.Now(), accountID)
	} else {
		result, err = r.db.Exec(query, newBalace, time.Now(), accountID)
	}
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check the update result %w", err)

	}
	if rowsAffected == 0 {
		return fmt.Errorf("account not found %w", err)
	}
	return nil
}

func (r *AccountRepository) GetUpdateBalance(tx *sql.Tx, accountID int) (float64, error) {
	query := `
	SELECT balance 
	FROM accounts
	WHERE id = $1
	FOR UPDATE
	`
	var balance float64

	err := tx.QueryRow(query, accountID).Scan(&balance)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("No account found")
	}

	if err != nil {
		return 0, fmt.Errorf("failed to get balance %w", err)
	}

	return balance, nil
}

func (r *AccountRepository) EmailExists(email string) (bool, error) {
	query := `
	SELECT EXISTS(
	SELECT 1 FROM accounts
	WHERE email = $1
	)
	`
	var exists bool

	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("account not found: %w", err)
	}
	return exists, nil
}

func (r *AccountRepository) Delete(id int) error {
	query := `
	INSERT accounts
	SET status = $1 , updated_at = $2
	WHERE id = $3
	`

	result, err := r.db.Exec(query, models.AccountStatusClosed, time.Now(), id)

	if err != nil {
		return fmt.Errorf("Failed to delete account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No account found")
	}

	return nil
}

func (r *AccountRepository) List(page, limit int) ([]*models.Account, int, error) {
	offset := (page - 1) * limit
	var totalCount int
	countQuery := `
	SELECT * FROM accounts WHERE status != $1
	`
	err := r.db.QueryRow(countQuery, models.AccountStatusClosed).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to get total count: %w", err)
	}
	query := `
	SELECT id, email, password_hash, first_name, last_name, balance, currency, status, created_at, updated_at
	FROM accounts
	WHERE status != $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, models.AccountStatusClosed, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list accounts: %w", err)
	}
	defer rows.Close()

	accounts := make([]*models.Account, 0)

	for rows.Next() {
		account := &models.Account{}
		err := rows.Scan(&account.ID, account.FirstName, account.LastName, account.Balance, account.Currency, account.Email, account.PasswordHash, account.Status, account.CreatedAt, account.UpdatedAt)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating accounts: %w", err)
	}
	return accounts, totalCount, nil
}
