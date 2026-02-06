package repository

import (
	"github.com/wizzyszn/go_bank/db"
)

type TransactionRepositoty struct {
	db *db.DB
}

func NewTransactionRepository(db *db.DB) *TransactionRepositoty {
	return &TransactionRepositoty{db: db}
}

func (t *TransactionRepositoty) Create(tx db.TxFunc) {

}
