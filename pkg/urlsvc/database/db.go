package database

import (
	"context"
	"database/sql"
	_ "embed"

	sqlc "github.com/cmgsj/urlshortener/pkg/gen/sqlc/url/v1"
	_ "github.com/mattn/go-sqlite3"
)

type DB interface {
	sqlc.Querier
	Ping(context.Context) error
	ExecTx(context.Context, TxFunc) error
	BeginTx(context.Context) (Tx, error)
}

type Tx interface {
	sqlc.Querier
	Rollback() error
	Commit() error
}

var _ DB = (*dbImpl)(nil)

type dbImpl struct {
	sqlc.Querier
	qrs *sqlc.Queries
	db  *sql.DB
}

func (db *dbImpl) Ping(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

type TxFunc func(sqlc.Querier) error

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
	return &txImpl{Querier: db.qrs.WithTx(tx), tx: tx}, nil
}

var _ Tx = (*txImpl)(nil)

type txImpl struct {
	sqlc.Querier
	tx *sql.Tx
}

func (tx *txImpl) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *txImpl) Commit() error {
	return tx.tx.Commit()
}
