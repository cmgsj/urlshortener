package db

import (
	"context"
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed sql/schema/schema.sql
	ddl string
)

func Connect(driverName, dataSourceName string) (DBTX, error) {
	return sql.Open(driverName, dataSourceName)
}

func Migrate(ctx context.Context, dbtx DBTX) (DBTX, error) {
	if _, err := dbtx.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}
	return dbtx, nil
}

func Must(dbtx DBTX, err error) DBTX {
	if err != nil {
		panic(err)
	}
	return dbtx
}

func MustPrepare(ctx context.Context, dbtx DBTX) Querier {
	q, err := Prepare(ctx, dbtx)
	if err != nil {
		panic(err)
	}
	return q
}
