// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createUrlStmt, err = db.PrepareContext(ctx, createUrl); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUrl: %w", err)
	}
	if q.deleteUrlStmt, err = db.PrepareContext(ctx, deleteUrl); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUrl: %w", err)
	}
	if q.getUrlStmt, err = db.PrepareContext(ctx, getUrl); err != nil {
		return nil, fmt.Errorf("error preparing query GetUrl: %w", err)
	}
	if q.listUrlsStmt, err = db.PrepareContext(ctx, listUrls); err != nil {
		return nil, fmt.Errorf("error preparing query ListUrls: %w", err)
	}
	if q.updateUrlStmt, err = db.PrepareContext(ctx, updateUrl); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUrl: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createUrlStmt != nil {
		if cerr := q.createUrlStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUrlStmt: %w", cerr)
		}
	}
	if q.deleteUrlStmt != nil {
		if cerr := q.deleteUrlStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUrlStmt: %w", cerr)
		}
	}
	if q.getUrlStmt != nil {
		if cerr := q.getUrlStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUrlStmt: %w", cerr)
		}
	}
	if q.listUrlsStmt != nil {
		if cerr := q.listUrlsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listUrlsStmt: %w", cerr)
		}
	}
	if q.updateUrlStmt != nil {
		if cerr := q.updateUrlStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUrlStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db            DBTX
	tx            *sql.Tx
	createUrlStmt *sql.Stmt
	deleteUrlStmt *sql.Stmt
	getUrlStmt    *sql.Stmt
	listUrlsStmt  *sql.Stmt
	updateUrlStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:            tx,
		tx:            tx,
		createUrlStmt: q.createUrlStmt,
		deleteUrlStmt: q.deleteUrlStmt,
		getUrlStmt:    q.getUrlStmt,
		listUrlsStmt:  q.listUrlsStmt,
		updateUrlStmt: q.updateUrlStmt,
	}
}