package database

import (
	"context"
	"database/sql"
)

type Transaction struct {
	tx   *sql.Tx
	done bool
}

func (t *Transaction) Commit() error {
	if t.done {
		return ErrTxDone
	}
	t.done = true
	return t.tx.Commit()
}

func (t *Transaction) Rollback() error {
	if t.done {
		return ErrTxDone
	}
	t.done = true
	return t.tx.Rollback()
}

func BeginTx(ctx context.Context, db *sql.DB, opts *sql.TxOptions) (*Transaction, error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Transaction{tx: tx}, nil
}
