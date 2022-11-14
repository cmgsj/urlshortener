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
	return db
}

func getUserByEmail(db *sql.DB, ctx context.Context, email string) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(ctx, "SELECT user_id, email, password FROM users WHERE email = ?", email).Scan(&user.UserId, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
