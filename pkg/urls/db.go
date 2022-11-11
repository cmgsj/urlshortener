package urls

import (
	"context"
	"database/sql"
	"log"
)

type UrlEntity struct {
	UrlId       string
	RedirectUrl string
}

func initSqliteDB(sqliteDbName string) *sql.DB {
	db, err := sql.Open("sqlite3", sqliteDbName)
	if err != nil {
		log.Fatal(err)
	}
	err = createTable(db, context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createTable(db *sql.DB, ctx context.Context) error {
	_, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS urls (url_id TEXT PRIMARY KEY, redirect_url TEXT UNIQUE)")
	return err
}

func getUrl(db *sql.DB, ctx context.Context, urlId string) (*UrlEntity, error) {
	u := &UrlEntity{}
	err := db.QueryRowContext(ctx, "SELECT * FROM urls WHERE url_id = ?", urlId).Scan(&u.UrlId, &u.RedirectUrl)
	return u, err
}

func createUrl(db *sql.DB, ctx context.Context, urlId string, redirectUrl string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO urls VALUES (?, ?)", urlId, redirectUrl)
	return err
}

func updateUrl(db *sql.DB, ctx context.Context, urlId string, redirectUrl string) error {
	_, err := db.ExecContext(ctx, "UPDATE urls SET redirect_url = ? WHERE url_id = ?", redirectUrl, urlId)
	return err
}

func deleteUrl(db *sql.DB, ctx context.Context, urlId string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM urls WHERE url_id = ?", urlId)
	return err
}
