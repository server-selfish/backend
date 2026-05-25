package pkg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}
type txManager struct {
	db     *pgxpool.Pool
	logger *zerolog.Logger
}

func NewTxManager(db *pgxpool.Pool, logger zerolog.Logger) TxManager {
	return &txManager{db: db, logger: &logger}
}

func (tm *txManager) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := tm.db.Begin(ctx)
	if err != nil {
		tm.logger.Error().Err(err).Msg("begin transaction failed")
		return err
	}

	committed := false
	defer func() {
		if !committed {
			if err := tx.Rollback(ctx); err != nil {
				tm.logger.Error().Err(err).Msg("failed to rollback transaction")
			}
		}
	}()

	if err := fn(tx); err != nil {
		tm.logger.Error().Err(err).Msg("fn inside transaction error")
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		tm.logger.Error().Err(err).Msg("commit transaction failed")
		return err
	}

	committed = true
	return nil
}
