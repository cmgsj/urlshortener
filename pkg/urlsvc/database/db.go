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
	Tx    *sql.DB
	Query sqlc.Querier
}

func New(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	querier, err := sqlc.Prepare(context.Background(), db)
	if err != nil {
		return nil, err
	}
	return &DB{Tx: db, Query: querier}, nil
}

func Migrate(db *DB) (*DB, error) {
	if _, err := db.Tx.ExecContext(context.Background(), ddl); err != nil {
		return nil, err
	}
	return db, nil
}

func Must(db *DB, err error) *DB {
	if err != nil {
		panic(err)
	}
	return db
}
