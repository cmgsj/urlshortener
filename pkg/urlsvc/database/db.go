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

type (
	DB struct {
		sqlc.Querier
		q  *sqlc.Queries
		db *sql.DB
	}
	Options struct {
		Driver  string
		URI     string
		Migrate bool
	}
)

func New(opts Options) (*DB, error) {
	db, err := sql.Open(opts.Driver, opts.URI)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	if opts.Migrate {
		if _, err = db.ExecContext(ctx, ddl); err != nil {
			return nil, err
		}
	}
	queries, err := sqlc.Prepare(ctx, db)
	if err != nil {
		return nil, err
	}
	return &DB{Querier: queries, q: queries, db: db}, nil
}

func Must(db *DB, err error) *DB {
	if err != nil {
		panic(err)
	}
	return db
}

type TxFunc func(sqlc.Querier) error

func (db *DB) ExecTx(ctx context.Context, txFunc TxFunc) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = txFunc(db.q.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) Ping(ctx context.Context) error {
	return db.db.PingContext(ctx)
}
