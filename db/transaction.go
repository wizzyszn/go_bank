package db

import (
	"context"
	"database/sql"
	"fmt"
)

type TxFunc func(*sql.Tx) error

func (db *DB) WithTransaction(ctx context.Context, fn TxFunc) error {

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
	
}

func (db *DB) ExecutionInTransaction(ctx context.Context, query string, args ...any) error {
	return db.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, query, args...)
		return err
	})
}
