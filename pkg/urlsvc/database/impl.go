package database

import (
	"context"
	"database/sql"
	_ "embed"

	sqlc "github.com/cmgsj/urlshortener/pkg/gen/sqlc/url/v1"
	_ "github.com/mattn/go-sqlite3"
)

var _ DB = (*dbImpl)(nil)

type dbImpl struct {
	sqlc.Querier
	qrs *sqlc.Queries
	db  *sql.DB
}

func newDBImpl(q *sqlc.Queries, db *sql.DB) *dbImpl {
	return &dbImpl{Querier: q, qrs: q, db: db}
}

func (db *dbImpl) Ping(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

func (db *dbImpl) ExecTx(ctx context.Context, txFunc TxFunc) error {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = txFunc(tx); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *dbImpl) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return newTxImpl(db.qrs.WithTx(tx), tx), nil
}

var _ Tx = (*txImpl)(nil)

type txImpl struct {
	sqlc.Querier
	tx *sql.Tx
}

func newTxImpl(q *sqlc.Queries, tx *sql.Tx) *txImpl {
	return &txImpl{Querier: q, tx: tx}
}

func (tx *txImpl) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *txImpl) Commit() error {
	return tx.tx.Commit()
}
