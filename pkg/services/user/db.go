package auth

import (
	"context"
	"database/sql"
	"urlshortener/pkg/logger"
)

type UserEntity struct {
	UserId   string
	Email    string
	Password string
}

func initSqliteDB(sqliteDbName string) *sql.DB {
	db, err := sql.Open("sqlite3", sqliteDbName)
	if err != nil {
		logger.Fatal("failed to open sqlite db:", err)
	}
	err = createUsersTable(db, context.Background())
	if err != nil {
		logger.Fatal("failed to create table:", err)
	}
	return db
}

func createUsersTable(db *sql.DB, ctx context.Context) error {
	_, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS users (user_id SERIAL PRIMARY KEY, email TEXT UNIQUE, password TEXT)")
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS email_index ON users (email)")
	return err
}

func getUserById(db *sql.DB, ctx context.Context, userId string) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(ctx, "SELECT user_id, email, password FROM users WHERE user_id = ?", userId).Scan(&user.UserId, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func getUserByEmail(db *sql.DB, ctx context.Context, email string) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(ctx, "SELECT user_id, email, password FROM users WHERE email = ?", email).Scan(&user.UserId, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func createUser(db *sql.DB, ctx context.Context, email string, password string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO users (email, password) VALUES (?, ?)", email, password)
	return err
}
