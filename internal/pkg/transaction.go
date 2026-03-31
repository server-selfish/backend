package pkg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/server-selfish/backend/internal/domain/repository"
	storage_infra "github.com/server-selfish/backend/internal/infra/storage"
)

type (
	TxRunner interface {
		RunInTx(ctx context.Context, fn func(q repository.Querier) error) error
	}
	txRunner struct {
		pool *pgxpool.Pool
	}
)

func NewTxRunner(pool *pgxpool.Pool) TxRunner {
	return &txRunner{pool: pool}
}

func (t *txRunner) RunInTx(
	ctx context.Context,
	fn func(q repository.Querier) error,
) error {
	pgxTx, err := t.pool.Begin(ctx)
	if err != nil {
		return err
	}

	tx := storage_infra.NewTx(pgxTx)

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
