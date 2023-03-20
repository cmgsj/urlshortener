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
	DB interface {
		sqlc.Querier
		Ping(context.Context) error
		ExecTx(context.Context, TxFunc) error
		BeginTx(context.Context) (Tx, error)
	}

	Tx interface {
		sqlc.Querier
		Rollback() error
		Commit() error
	}

	TxFunc func(sqlc.Querier) error

	Options struct {
		Driver      string
		ConnString  string
		AutoMigrate bool
	}
)

func New(opts Options) (DB, error) {
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
	return newDBImpl(q, db), nil
}

func Must(db DB, err error) DB {
	if err != nil {
		panic(err)
	}
	return db
}
