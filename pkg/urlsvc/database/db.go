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
		Tx    *sql.DB
		Query sqlc.Querier
	}
	Options struct {
		Driver  string
		URI     string
		Migrate bool
	}
)

func New(opt Options) (*DB, error) {
	tx, err := sql.Open(opt.Driver, opt.URI)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	if opt.Migrate {
		if _, err = tx.ExecContext(ctx, ddl); err != nil {
			return nil, err
		}
	}
	querier, err := sqlc.Prepare(ctx, tx)
	if err != nil {
		return nil, err
	}
	return &DB{Tx: tx, Query: querier}, nil
}

func Must(db *DB, err error) *DB {
	if err != nil {
		panic(err)
	}
	return db
}

func DDL() string {
	return ddl
}
