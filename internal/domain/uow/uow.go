package uow

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgconn"
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

	ctx = context.WithValue(ctx, TxKey{}, tx)

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
    Query(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
    QueryRow(ctx context.Context, sql string, args ...any) *sql.Row
}

func ExtractExecutor(ctx context.Context, db *sql.DB) Executor {
	tx, ok := ctx.Value(TxKey{}).(*sql.Tx)
	if !ok {
		return db
	}
	return tx
}
