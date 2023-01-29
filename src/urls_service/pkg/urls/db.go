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
	err := db.QueryRowContext(ctx, "SELECT * FROM urls WHERE url_id = ?", urlId).Scan(&u.UrlId, &u.RedirectUrl)
	return u, err
}

func createUrl(ctx context.Context, db *sql.DB, urlId string, redirectUrl string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO urls (url_id, redirect_url) VALUES (?, ?)", urlId, redirectUrl)
	return err
}

func updateUrl(ctx context.Context, db *sql.DB, urlId string, redirectUrl string) error {
	_, err := db.ExecContext(ctx, "UPDATE urls SET redirect_url = ? WHERE url_id = ?", redirectUrl, urlId)
	return err
}

func deleteUrl(ctx context.Context, db *sql.DB, urlId string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM urls WHERE url_id = ?", urlId)
	return err
}
