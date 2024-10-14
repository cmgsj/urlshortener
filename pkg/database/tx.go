package database

import (
	"context"
	"database/sql"
)

type Txer interface {
	Tx(ctx context.Context, tx *sql.Tx) error
}

type TxFunc func(ctx context.Context, tx *sql.Tx) error

func (f TxFunc) Tx(ctx context.Context, tx *sql.Tx) error {
	return f(ctx, tx)
}

func Tx(ctx context.Context, db *sql.DB, opts *sql.TxOptions, t Txer) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	err = t.Tx(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
