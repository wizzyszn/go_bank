package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/wizzyszn/go_bank/db"
	"github.com/wizzyszn/go_bank/models"
	"github.com/wizzyszn/go_bank/repository"
	"github.com/wizzyszn/go_bank/utils"
)

type TransactionService struct {
	db              *db.DB
	accountRepo     *repository.AccountRepository
	transactionRepo *repository.TransactionRepositoty
}

func NewTransactionService(
	database *db.DB,
	accountRepo *repository.AccountRepository,
	transactionRepo *repository.TransactionRepositoty,
) *TransactionService {
	return &TransactionService{
		db:              database,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *TransactionService) Deposit(accountID int, req *models.DepositRequest) (*models.TransactionResponse, error) {

	if err := utils.ValidateAmount(req.Amount); err != nil {
		return nil, err
	}

	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	if account.Status != models.AccountStatusActice {
		return nil, fmt.Errorf("account is %s", account.Status)
	}

	var transaction *models.Transaction

	err = s.db.WithTransaction(context.Background(), func(tx *sql.Tx) error {
		currentBalance, err := s.accountRepo.GetBalanceForUpdate(tx, accountID)
		if err != nil {
			return err
		}

		newBalance := currentBalance + req.Amount

		if err := s.accountRepo.UpdateBalance(tx, accountID, newBalance); err != nil {
			return err
		}

		transaction, err = s.transactionRepo.Create(tx, nil, &accountID, req.Amount, models.TransactionTypeDeposit, req.Description)

		return err
	})

	if err != nil {
		return nil, fmt.Errorf("deposit failed: %w", err)
	}
	return transaction.ToResponse(), nil
}

func (s *TransactionService) WithDraw(accountID int, req *models.WitdrawRequest) (*models.TransactionResponse, error) {
	if err := utils.ValidateAmount(req.Amount); err != nil {
		return nil, err
	}
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}
	if account.Status != models.AccountStatusActice {
		return nil, fmt.Errorf("account is %s", account.Status)
	}

	var transaction *models.Transaction

	err = s.db.WithTransaction(context.Background(), func(tx *sql.Tx) error {
		currentBalance, err := s.accountRepo.GetBalanceForUpdate(tx, accountID)
		if err != nil {
			return err
		}
		if currentBalance < req.Amount {
			return fmt.Errorf("insufficient funds: have %.2f, need %.2f", currentBalance, req.Amount)
		}
		newBalace := currentBalance - req.Amount

		err = s.accountRepo.UpdateBalance(tx, accountID, newBalace)
		if err != nil {
			return err
		}

		transaction, err = s.transactionRepo.Create(tx, &accountID, nil, req.Amount, models.TransactionTypeWithdraw, req.Description)

		return err
	})

	if err != nil {
		return nil, fmt.Errorf("withdrawal failed: %w", err)
	}

	return transaction.ToResponse(), nil
}

func (s *TransactionService) Transfer(fromAccountID int, req models.TransferRequest) (*models.TransactionResponse, error) {
	if err := utils.ValidateAmount(req.Amount); err != nil {
		return nil, err
	}

	if err := utils.ValidateAccountID(req.ToAccountID); err != nil {
		return nil, err
	}

	if fromAccountID == req.ToAccountID {
		return nil, fmt.Errorf("cannot transfer to your own account")
	}

	fromAccount, err := s.accountRepo.GetByID(fromAccountID)
	if err != nil {
		return nil, fmt.Errorf("sender account not found")
	}
	if fromAccount.Status != models.AccountStatusActice {
		return nil, fmt.Errorf("sender account is %s", fromAccount.Status)
	}

	toAccount, err := s.accountRepo.GetByID(req.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("recipient account not found")
	}

	if toAccount.Status != models.AccountStatusActice {
		return nil, fmt.Errorf("sender account is %s", fromAccount.Status)
	}

	var transaction *models.Transaction

	err = s.db.WithTransaction(context.Background(), func(tx *sql.Tx) error {
		firstID, secondID := fromAccountID, req.ToAccountID

		if firstID > secondID {
			firstID, secondID = secondID, firstID
		}
		firstBalance, err := s.accountRepo.GetBalanceForUpdate(tx, firstID)
		if err != nil {
			return err
		}
		secondBalance, err := s.accountRepo.GetBalanceForUpdate(tx, secondID)
		if err != nil {
			return err
		}

		senderBalance := firstBalance
		receiverBalance := secondBalance
		if fromAccountID > req.ToAccountID {
			senderBalance, receiverBalance = secondBalance, firstBalance
		}
		if senderBalance < req.Amount {
			return fmt.Errorf("insufficient funds: have %.2f, need %.2f", senderBalance, req.Amount)
		}
		if err := s.accountRepo.UpdateBalance(tx, fromAccountID, senderBalance-req.Amount); err != nil {
			return err
		}
		if err := s.accountRepo.UpdateBalance(tx, req.ToAccountID, receiverBalance+req.Amount); err != nil {
			return err
		}
		toAccountID := req.ToAccountID
		transaction, err = s.transactionRepo.Create(
			tx,
			&fromAccountID,
			&toAccountID,
			req.Amount,
			models.TransactionTypeTransfer,
			req.Description,
		)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("transfer failed: %w", err)
	}

	return transaction.ToResponse(), nil
}
