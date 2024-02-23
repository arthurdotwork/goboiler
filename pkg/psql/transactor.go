package psql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
)

type Transactor interface {
	WithinTransaction(context.Context, func(ctx context.Context) error) error
}

type transactor struct {
	db Queryable
}

func NewTransactor(db Queryable) Transactor {
	return &transactor{db: db}
}

type txCtxKey struct{}

func txToContext(ctx context.Context, tx Queryable) context.Context {
	return context.WithValue(ctx, txCtxKey{}, &TxQueryable{Queryable: tx})
}

func TxFromContext(ctx context.Context) Queryable {
	tx, ok := ctx.Value(txCtxKey{}).(Queryable)
	if ok {
		return tx
	}

	return nil
}

func (t *transactor) WithinTransaction(ctx context.Context, txFunc func(txCtx context.Context) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txCtx := txToContext(ctx, tx)
	if err := txFunc(txCtx); err != nil {
		tx.Rollback() // nolint: errcheck // If rollback fails, there's nothing to do, the transaction will expire by itself
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

type TxQueryable struct {
	Queryable
	sync.Mutex
}

func (t *TxQueryable) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	t.Lock()
	defer t.Unlock()

	return t.Queryable.ExecContext(ctx, query, args...) // nolint:wrapcheck
}

func (t *TxQueryable) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	t.Lock()
	defer t.Unlock()

	return t.Queryable.GetContext(ctx, dest, query, args...) // nolint:wrapcheck
}

func (t *TxQueryable) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	t.Lock()
	defer t.Unlock()

	return t.Queryable.SelectContext(ctx, dest, query, args...) // nolint:wrapcheck
}

func (t *TxQueryable) NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error) {
	t.Lock()
	defer t.Unlock()

	return t.Queryable.NamedExecContext(ctx, query, arg) // nolint:wrapcheck
}

func (t *TxQueryable) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	t.Lock()
	defer t.Unlock()

	return t.Queryable.PrepareNamedContext(ctx, query) // nolint:wrapcheck
}

func (t *TxQueryable) PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	t.Lock()
	defer t.Unlock()

	return t.Queryable.PreparexContext(ctx, query) // nolint:wrapcheck
}
