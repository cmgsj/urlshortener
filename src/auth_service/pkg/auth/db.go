package auth

import (
	"auth_service/pkg/database"
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type UserEntity struct {
	UserId   int64
	Email    string
	Password string
}

func (s *Service) intiDB(sqliteDbName string) {
	db, err := sql.Open("sqlite3", sqliteDbName)
	if err != nil {
		s.logger.Fatal("failed to open sqlite db:", zap.Error(err))
	}
	err = database.CreateTables(context.Background(), db)
	if err != nil {
		s.logger.Fatal("failed to create tables:", zap.Error(err))
	}
	err = database.SeedDB(context.Background(), db)
	if err != nil {
		s.logger.Error("failed to seed db:", zap.Error(err))
	} else {
		s.logger.Info("db seeded")
	}
	s.db = db
}

func getUserByEmail(ctx context.Context, db *sql.DB, email string) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(ctx, "SELECT user_id, email, password FROM users WHERE email = ?", email).Scan(&user.UserId, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func createUser(ctx context.Context, db *sql.DB, email, password string) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(ctx, "INSERT INTO users (email, password) VALUES (?, ?)", email, password).Scan(&user.UserId, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
