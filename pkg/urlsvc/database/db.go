package database

import (
	"context"
	"database/sql"
	_ "embed"

	sqlc "github.com/cmgsj/urlshortener/pkg/gen/sqlc/url/v1"
	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed sql/v1/schema/schema.sql
	ddl string
)

type DB struct {
	Query sqlc.Querier
	qrs   *sqlc.Queries
	db    *sql.DB
}

type Options struct {
	Driver      string
	ConnString  string
	AutoMigrate bool
}

func New(opts Options) (*DB, error) {
	db, err := sql.Open(opts.Driver, opts.ConnString)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	if opts.AutoMigrate {
		if _, err = db.ExecContext(ctx, ddl); err != nil {
			return nil, err
		}
	}
	q, err := sqlc.Prepare(ctx, db)
	if err != nil {
		return nil, err
	}
	return &DB{Query: q, qrs: q, db: db}, nil
}

func Must(db *DB, err error) *DB {
	if err != nil {
		panic(err)
	}
	return db
}

func (db *DB) Ping(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

type TxFunc func(sqlc.Querier) error

func (db *DB) ExecTx(ctx context.Context, txFunc TxFunc) error {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = txFunc(tx.Query); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{Query: db.qrs.WithTx(tx), tx: tx}, nil
}

type Tx struct {
	Query sqlc.Querier
	tx    *sql.Tx
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}
