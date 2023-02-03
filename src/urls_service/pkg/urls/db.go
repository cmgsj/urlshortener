package urls

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

type UrlEntity struct {
	UrlId       string
	RedirectUrl string
}

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

func createTables(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, CreateTablesSQL)
	return err
}

func seedDB(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, SeedDBSQL)
	return err
}

func (s *Service) IntiDB(sqliteDbName string) {
	db, err := sql.Open("sqlite3", sqliteDbName)
	fmt.Println("sqliteDbName: ", sqliteDbName)
	if err != nil {
		s.Logger.Fatal("failed to open sqlite db:", zap.Error(err))
	}
	err = createTables(context.Background(), db)
	if err != nil {
		s.Logger.Fatal("failed to create tables:", zap.Error(err))
	}
	err = seedDB(context.Background(), db)
	if err != nil {
		s.Logger.Error("failed to seed db:", zap.Error(err))
	} else {
		s.Logger.Info("db seeded")
	}
	s.Db = db
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
