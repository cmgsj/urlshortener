package database

import (
	"context"
	"database/sql"
	_ "embed"

	sqlc "github.com/cmgsj/urlshortener/pkg/gen/sqlc/urls/v1"
	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed sql/urls/v1/schema.sql
	urlsDDL string
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

type TxFunc func(sqlc.Querier) error

type Options struct {
	Driver      string
	ConnString  string
	AutoMigrate bool
}

func New(o Options) (DB, error) {
	db, err := sql.Open(o.Driver, o.ConnString)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	if o.AutoMigrate {
		if _, err = db.ExecContext(ctx, urlsDDL); err != nil {
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
