package pkg

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}
type txManager struct {
	db *pgxpool.Pool
}

func NewTxManager(db *pgxpool.Pool) TxManager {
	return &txManager{db: db}
}

func (tm *txManager) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := tm.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {

		if err := tx.Rollback(ctx); err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
