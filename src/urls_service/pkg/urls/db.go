package urls

import (
	"context"
	"database/sql"
	"urls_service/pkg/database"

	"go.uber.org/zap"
)

type UrlEntity struct {
	UrlId       string
	RedirectUrl string
	UserId      int64
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

func getUrlById(ctx context.Context, db *sql.DB, urlId string) (*UrlEntity, error) {
	u := &UrlEntity{}
	err := db.QueryRowContext(ctx, "SELECT * FROM urls WHERE url_id = ?", urlId).Scan(&u.UrlId, &u.RedirectUrl, &u.UserId)
	return u, err
}

// func getUrlsByUserId(ctx context.Context, db *sql.DB, userId int64) ([]*UrlEntity, error) {
// 	rows, err := db.QueryContext(ctx, "SELECT * FROM urls WHERE user_id = ?", userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var urls []*UrlEntity
// 	for rows.Next() {
// 		u := &UrlEntity{}
// 		err := rows.Scan(&u.UrlId, &u.RedirectUrl, &u.UserId)
// 		if err != nil {
// 			return nil, err
// 		}
// 		urls = append(urls, u)
// 	}
// 	return urls, nil
// }

func createUrl(ctx context.Context, db *sql.DB, urlId string, redirectUrl string, userId int64) error {
	_, err := db.ExecContext(ctx, "INSERT INTO urls (url_id, redirect_url, user_id) VALUES (?, ?, ?)", urlId, redirectUrl, userId)
	return err
}

func updateUrl(ctx context.Context, db *sql.DB, urlId string, redirectUrl string, userId int64) error {
	_, err := db.ExecContext(ctx, "UPDATE urls SET redirect_url = ? WHERE url_id = ? AND user_id = ?", redirectUrl, urlId, userId)
	return err
}

func deleteUrl(ctx context.Context, db *sql.DB, urlId string, userId int64) error {
	_, err := db.ExecContext(ctx, "DELETE FROM urls WHERE url_id = ? AND user_id = ?", urlId, userId)
	return err
}
