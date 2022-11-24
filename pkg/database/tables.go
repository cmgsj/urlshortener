package database

import (
	"context"
	"database/sql"
)

var (
	CreateTablesSQL = `
		CREATE TABLE IF NOT EXISTS users
		(
			user_id  SERIAL PRIMARY KEY,
			email    TEXT UNIQUE,
			password TEXT
		);
		CREATE INDEX IF NOT EXISTS email_index ON users (email);
		CREATE TABLE IF NOT EXISTS urls
		(
			url_id       TEXT PRIMARY KEY,
			redirect_url TEXT UNIQUE,
			user_id      INTEGER,
			FOREIGN KEY (user_id) REFERENCES users (user_id)
		);
		CREATE INDEX IF NOT EXISTS redirect_url_index ON urls (redirect_url);
	`
	SeedDBSQL = `
		INSERT INTO users
		VALUES (1, 'user1@example.com', 'user1'),
			   (2, 'user2@example.com', 'user2'),
			   (3, 'user3@example.com', 'user3');
		
		INSERT INTO urls
		VALUES ('abcdef01', 'https://www.google.com', 1),
			   ('abcdef02', 'https://www.youtube.com', 2),
			   ('abcdef03', 'https://www.apple.com', 3);
		`
)

func CreateTables(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, CreateTablesSQL)
	return err
}

func SeedDB(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, SeedDBSQL)
	return err
}
