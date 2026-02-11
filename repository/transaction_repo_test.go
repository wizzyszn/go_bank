package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wizzyszn/go_bank/db"
	"github.com/wizzyszn/go_bank/models"
)

func setupTestDB(t *testing.T) (*db.DB, func()) {
	// Try to load .env from root
	wd, _ := os.Getwd()
	parent := filepath.Dir(wd)
	err := godotenv.Load(filepath.Join(parent, ".env"))
	if err != nil {
		// Try loading from current dir
		_ = godotenv.Load(".env")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	// valid default if env vars missing (based on seen .env)
	if os.Getenv("DB_HOST") == "" {
		dsn = "host=localhost port=5432 user=postgres password=yourpassword dbname=bankdb sslmode=disable"
	}

	cfg := db.Config{
		DSN:             dsn,
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}

	database, err := db.New(cfg)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	return database, func() {
		database.Close()
	}
}

func TestGetTotalBalance(t *testing.T) {
	database, teardown := setupTestDB(t)
	defer teardown()

	accountRepo := NewAccountRepository(database)
	transactionRepo := NewTransactionRepository(database)

	// Create test accounts
	// Email must be unique, so use random or timestamp
	timestamp := time.Now().UnixNano()
	acc1, err := accountRepo.Create(fmt.Sprintf("test1_%d@example.com", timestamp), "hash", "John", "Doe")
	if err != nil {
		t.Fatalf("failed to create account 1: %v", err)
	}
	acc2, err := accountRepo.Create(fmt.Sprintf("test2_%d@example.com", timestamp), "hash", "Jane", "Doe")
	if err != nil {
		t.Fatalf("failed to create account 2: %v", err)
	}

	// Cleanup accounts after test
	defer func() {
		// We don't have a Delete method in AccountRepo that creates hard delete,
		// and schema has cascades?
		// accounts.id is referenced by transactions.
		// Let's rely on random emails to avoid collision and maybe manual cleanup if needed.
		// For local dev DB, we might want to clean up.
		// models.AccountStatusClosed logic in Delete method updates status, doesn't remove row.
		// To keep test clean, we could delete manually.
		_, _ = database.Exec("DELETE FROM transactions WHERE from_account_id = $1 OR to_account_id = $1", acc1.ID)
		_, _ = database.Exec("DELETE FROM transactions WHERE from_account_id = $1 OR to_account_id = $1", acc2.ID)
		_, _ = database.Exec("DELETE FROM accounts WHERE id = $1", acc1.ID)
		_, _ = database.Exec("DELETE FROM accounts WHERE id = $1", acc2.ID)
	}()

	// 1. Initial Balance should be 0
	bal, err := transactionRepo.GetTotalBalance(acc1.ID)
	if err != nil {
		t.Fatalf("GetTotalBalance failed: %v", err)
	}
	if bal != 0 {
		t.Errorf("expected initial balance 0, got %f", bal)
	}

	// 2. Deposit 1000 to Acc1
	// For deposit, From can be same as To or nil?
	// Based on implementation of GetTotalBalance:
	// + (to_account_id == id)
	// - (from_account_id == id)
	// If we set from = to for deposit, it cancels out?
	// Check create transaction logic:
	// Create(..., fromAccountID, toAccountID, ...)
	// If Deposit: usually money comes from "outside".
	// Schema allows from_account_id to be NULL?
	// CHECK (from_account_id IS NOT NULL OR to_account_id IS NOT NULL)
	// So for deposit, likely from_account_id is NULL or a system account.
	// Let's try passing something else or creating a system account.
	// Or maybe check how Create handles it. Create takes *int.

	// Let's create a "Bank" account or use nil for fromAccountID
	// But Create takes *int.
	// Let's try passing nil for fromAccountID.

	// Acc1 Deposit 1000
	amount1000 := 1000.00
	_, err = transactionRepo.Create(nil, nil, &acc1.ID, amount1000, models.TransactionTypeDeposit, "Deposit")
	if err != nil {
		t.Fatalf("failed to create deposit: %v", err)
	}

	bal, err = transactionRepo.GetTotalBalance(acc1.ID)
	if err != nil {
		t.Fatalf("GetTotalBalance failed: %v", err)
	}
	if bal != 1000 {
		t.Errorf("expected balance 1000 after deposit, got %f", bal)
	}

	// 3. Transfer 200 from Acc1 to Acc2
	amount200 := 200.00
	_, err = transactionRepo.Create(nil, &acc1.ID, &acc2.ID, amount200, models.TransactionTypeTransfer, "Transfer to Acc2")
	if err != nil {
		t.Fatalf("failed to create transfer: %v", err)
	}

	// Acc1 Balance should be 800
	bal, err = transactionRepo.GetTotalBalance(acc1.ID)
	if err != nil {
		t.Fatalf("GetTotalBalance failed: %v", err)
	}
	if bal != 800 {
		t.Errorf("expected balance 800 after transfer, got %f", bal)
	}

	// Acc2 Balance should be 200
	bal2, err := transactionRepo.GetTotalBalance(acc2.ID)
	if err != nil {
		t.Fatalf("GetTotalBalance failed for acc2: %v", err)
	}
	if bal2 != 200 {
		t.Errorf("expected balance 200 for acc2, got %f", bal2)
	}

	// 4. Withdraw 100 from Acc1
	amount100 := 100.00
	// Withdraw: From Acc1, To nil?
	_, err = transactionRepo.Create(nil, &acc1.ID, nil, amount100, models.TransactionTypeWithdraw, "Withdrawal")
	if err != nil {
		t.Fatalf("failed to create withdrawal: %v", err)
	}

	// Acc1 Balance should be 700
	bal, err = transactionRepo.GetTotalBalance(acc1.ID)
	if err != nil {
		t.Fatalf("GetTotalBalance failed: %v", err)
	}
	if bal != 700 {
		t.Errorf("expected balance 700 after withdrawal, got %f", bal)
	}
}
