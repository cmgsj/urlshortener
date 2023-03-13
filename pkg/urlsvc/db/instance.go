package db

import (
	"context"
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db Querier
	//go:embed sql/schema/schema.sql
	ddl string
)

func Instance() Querier {
	return db
}

func SetInstance(dbtx DBTX) Querier {
	db = New(dbtx)
	return Instance()
}

func Connect(driverName, dataSourceName string) (DBTX, error) {
	dbtx, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return dbtx, nil
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
