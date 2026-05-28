package uow

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UoW interface {
	Do(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
}

type SQLUoW struct {
	db *sql.DB
}

func NewSQLUoW(db *sql.DB) *SQLUoW {
	return &SQLUoW{db: db}
}

type txKey struct{}

func (u *SQLUoW) Do(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error) {
	isCommited := false

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if !isCommited {
			_ = tx.Rollback()
		}
	}()

	ctx = context.WithValue(ctx, txKey{}, tx)

	result, err := fn(ctx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	isCommited = true

	return result, nil
}

type Executor interface {
    Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
    Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
    QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func GetExecutor(ctx context.Context, db *pgxpool.Pool) Executor {
    if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
        return tx
    }
    return db
}
