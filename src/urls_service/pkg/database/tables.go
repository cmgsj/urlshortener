package database

import (
	"context"
	"database/sql"
)

var (
	CreateTablesSQL = `
		CREATE TABLE IF NOT EXISTS urls
		(
			url_id       TEXT PRIMARY KEY,
			redirect_url TEXT UNIQUE
		);
		CREATE INDEX IF NOT EXISTS redirect_url_index ON urls (redirect_url);`
	SeedDBSQL = `
		INSERT INTO urls
		VALUES ('abcdef01', 'https://www.google.com'),
			   ('abcdef02', 'https://www.youtube.com'),
			   ('abcdef03', 'https://www.apple.com');`
)

func CreateTables(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, CreateTablesSQL)
	return err
}

func SeedDB(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, SeedDBSQL)
	return err
}
