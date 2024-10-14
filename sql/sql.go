package sql

import (
	"context"
	"database/sql"
	_ "embed"
)

var (
	//go:embed schema.sql
	schema string
)

func Schema() string {
	return schema
}

func Migrate(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, schema)
	return err
}
