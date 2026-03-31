package storage_infra

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/server-selfish/backend/internal/domain/repository"
)

type (
	txImpl struct {
		tx pgx.Tx
	}
	querier struct {
		pgx *pgxpool.Pool
	}
)

func NewQuerier(pool *pgxpool.Pool) repository.Querier {
	return &querier{pgx: pool}
}

func (p *querier) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pgx.QueryRow(ctx, sql, args...)
}
func (p *querier) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return p.pgx.Exec(ctx, sql, args...)
}

func (p *querier) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pgx.Query(ctx, sql, args...)
}

func NewTx(tx pgx.Tx) repository.Tx {
	return &txImpl{tx: tx}
}

func (t *txImpl) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *txImpl) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

// Exec implements [Tx].
func (t *txImpl) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return t.Exec(ctx, sql, args...)
}

// Query implements [Tx].
func (t *txImpl) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return t.Query(ctx, sql, args...)
}

// QueryRow implements [Tx].
func (t *txImpl) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return t.QueryRow(ctx, sql, args...)
}
