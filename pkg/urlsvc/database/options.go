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

type Options struct {
	Driver      string
	ConnString  string
	AutoMigrate bool
}

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
	return &dbImpl{Querier: q, qrs: q, db: db}, nil
}

func Must(db DB, err error) DB {
	if err != nil {
		panic(err)
	}
	return db
}
