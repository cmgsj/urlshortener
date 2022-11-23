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

func intiDB(sqliteDbName string) *sql.DB {
	db, err := sql.Open("sqlite3", sqliteDbName)
	if err != nil {
		logger.Fatal("failed to open sqlite db:", err)
	}
	return db
}

func getUserById(ctx context.Context, db *sql.DB, userId string) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(ctx, "SELECT user_id, email, password FROM users WHERE user_id = ?", userId).Scan(&user.UserId, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func getUserByEmail(ctx context.Context, db *sql.DB, email string) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(ctx, "SELECT user_id, email, password FROM users WHERE email = ?", email).Scan(&user.UserId, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func createUser(ctx context.Context, db *sql.DB, email string, password string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO users (email, password) VALUES (?, ?)", email, password)
	return err
}
