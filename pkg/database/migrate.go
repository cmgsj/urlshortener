package database

import (
	"context"
	"database/sql"
	_ "embed"
)

var (
	//go:embed sql/schema.sql
	ddl string
)

func DDL() string {
	return ddl
}

func Migrate(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, ddl)
	return err
}
